package repository

import (
	"basicLoginRest/config"
	"basicLoginRest/internal/session"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

var (
	jsonMarshal                      = json.Marshal
	jsonUnmarshal                    = json.Unmarshal
	_             session.Repository = &cacheManager{}
	_             session.Store      = &store{}
)

func NewCacheManager(cli *redis.Client, cfg *config.Config, prefix ...string) session.Repository {
	cache := &cacheManager{
		cli: cli,
		cfg: cfg,
	}

	if len(prefix) > 0 {
		cache.prefix = prefix[0]
	}

	return cache
}

type cacheManager struct {
	cli *redis.Client
	cfg *config.Config
	prefix string
}

func (cm *cacheManager) getKey(sid string) string {
	return cm.prefix + sid
}

func (cm *cacheManager) getValue(ctx context.Context, sid string) (string, error) {
	res, err := cm.cli.Get(ctx, cm.getKey(sid)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return res, nil
}

func (cm *cacheManager) parseVal(value string) (map[string]interface{}, error) {
	var values map[string]interface{}
	if len(value) > 0 {
		if err := jsonUnmarshal([]byte(value), &values); err != nil {
			return nil, err
		}
	}
	return values, nil
}

func (cm *cacheManager) Check(ctx context.Context, sid string) (bool, error) {
	res, err := cm.cli.Exists(ctx, cm.getKey(sid)).Result()
	//1 - exists 0 - does not exist
	return res == 1, err
}

func (cm *cacheManager) Create(ctx context.Context, sid string, expired time.Duration) (session.Store, error) {
	return NewStore(ctx, cm, sid, expired, nil), nil
}

func (cm *cacheManager) Delete(ctx context.Context, sid string) error {
	if ok, err := cm.Check(ctx, sid); err != nil {
		return err
	} else if !ok {
		return nil
	}
	return cm.cli.Del(ctx, cm.getKey(sid)).Err()
}

func (cm *cacheManager) Update(ctx context.Context, sid string, expired time.Duration) (session.Store, error) {
	value, err := cm.getValue(ctx, sid)
	if err != nil {
		return nil, err
	} else if len(value) == 0 {
		return NewStore(ctx, cm, sid, expired,nil), nil
	}

	if err := cm.cli.Expire(ctx, sid, expired).Err(); err != nil {
		return nil, err
	}

	values, err := cm.parseVal(value)
	if err != nil {
		return nil, err
	}

	return NewStore(ctx, cm, sid, expired, values), nil
}

func (cm *cacheManager) Refresh(ctx context.Context, oldSid, newSid string, expired time.Duration) (session.Store, error) {
	value, err := cm.getValue(ctx, oldSid)
	if err != nil {
		return nil, err
	} else if len(value) == 0 {
		return NewStore(ctx, cm, newSid, expired,nil), nil
	}

	pipe := cm.cli.TxPipeline()
	pipe.Set(ctx, cm.getKey(newSid), value, expired)
	pipe.Del(ctx, cm.getKey(oldSid))

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	values, err := cm.parseVal(value)
	if err != nil {
		return nil, err
	}

	return NewStore(ctx, cm, newSid, expired, values), nil
}

func (cm *cacheManager) Connect() (bool, error) {
	res, err := cm.cli.Ping(context.TODO()).Result()
	return res == "PONG", err
}

func (cm *cacheManager) Close() error {
	return cm.cli.Close()
}


func NewStore(ctx context.Context, cm *cacheManager, sid string, expired time.Duration, values map[string]interface{}) session.Store {
	if values == nil {
		values = make(map[string]interface{})
	}

	return &store{
		ctx:     ctx,
		sid:     sid,
		cm:      cm,
		expired: expired,
		values:  values,
	}
}

type store struct {
	sync.RWMutex
	ctx context.Context
	sid string
	cm *cacheManager
	expired time.Duration
	values  map[string]interface{}
}

func (c *store) Context() context.Context {
	return c.ctx
}

func (c *store) SessionID() string {
	return c.sid
}

func (c *store) Set(key string, value interface{}) {
	c.Lock()
	c.values[key] = value
	c.Unlock()
}

func (c *store) Get(key string) (interface{}, bool) {
	c.RLock()
	val, ok := c.values[key]
	c.RUnlock()
	return val, ok
}

func (c *store) Delete(key string) interface{} {
	c.RLock()
	v, ok := c.values[key]
	c.RUnlock()
	if ok {
		c.Lock()
		delete(c.values, key)
		c.Unlock()
	}
	return v
}

func (c *store) Save() error {
	var buf []byte
	var err error

	c.RLock()
	if len(c.values) > 0 {
		buf, err = jsonMarshal(c.values)
		if err != nil {
			c.RUnlock()
			return err
		}
	}
	c.RUnlock()

	return c.cm.cli.Set(c.ctx, c.cm.getKey(c.sid), string(buf), c.expired).Err()
}

func (c *store) Flush() error {
	c.Lock()
	c.values = make(map[string]interface{})
	c.Unlock()
	return c.Save()
}
