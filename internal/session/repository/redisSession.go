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
	jsonMarshal                     = json.Marshal
	jsonUnmarshal                   = json.Unmarshal
	_             session.UCSession = &cacheManager{}
	_             session.Session   = &cache{}
)

func NewRedisCache(opts *config.Config, prefix ...string) session.UCSession {
	if opts == nil {
		panic("redis options not specified")
	}

	redisOpts := redis.Options{}
	redisOpts.Addr = opts.Redis.Addr
	redisOpts.DB = opts.Redis.Database
	redisOpts.Password = opts.Redis.Password
	redisOpts.MaxRetries = opts.Redis.MaxRetries
	redisOpts.MaxRetryBackoff = opts.Redis.MaxRetryBackoff

	return NewRedisCacheCli(
		redis.NewClient(&redisOpts),
		prefix...)
}

func NewRedisCacheCli(cli *redis.Client, prefix ...string) session.UCSession {
	cache := &cacheManager{
		cli: cli,
	}

	if len(prefix) > 0 {
		cache.prefix = prefix[0]
	}

	return cache
}

type cacheManager struct {
	cli *redis.Client
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

func (cm *cacheManager) Create(ctx context.Context, sid string, expired time.Duration) (session.Session, error) {
	return newCache(ctx, cm, sid, expired, nil), nil
}

func (cm *cacheManager) Update(ctx context.Context, sid string, expired time.Duration) (session.Session, error) {
	value, err := cm.getValue(ctx, sid)
	if err != nil {
		return nil, err
	} else if len(value) == 0 {
		return newCache(ctx, cm, sid, expired,nil), nil
	}

	if err := cm.cli.Expire(ctx, sid, expired).Err(); err != nil {
		return nil, err
	}

	values, err := cm.parseVal(value)
	if err != nil {
		return nil, err
	}

	return newCache(ctx, cm, sid, expired, values), nil
}

func (cm *cacheManager) Refresh(ctx context.Context, oldSid, newSid string, expired time.Duration) (session.Session, error) {
	value, err := cm.getValue(ctx, oldSid)
	if err != nil {
		return nil, err
	} else if len(value) == 0 {
		return newCache(ctx, cm, newSid, expired,nil), nil
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

	return newCache(ctx, cm, newSid, expired, values), nil
}

func (cm *cacheManager) Connect() (bool, error) {
	res, err := cm.cli.Ping(context.TODO()).Result()
	return res == "PONG", err
}

func (cm *cacheManager) Close() error {
	return cm.cli.Close()
}


func newCache(ctx context.Context, cm *cacheManager, sid string, expired time.Duration, values map[string]interface{}) *cache {
	if values == nil {
		values = make(map[string]interface{})
	}

	return &cache{
		ctx:     ctx,
		sid:     sid,
		cm:      cm,
		expired: expired,
		values:  values,
	}
}

type cache struct {
	sync.RWMutex
	ctx context.Context
	sid string
	cm *cacheManager
	expired time.Duration
	values  map[string]interface{}
}

func (c *cache) Context() context.Context {
	return c.ctx
}

func (c *cache) SessionID() string {
	return c.sid
}

func (c *cache) Set(key string, value interface{}) {
	c.Lock()
	c.values[key] = value
	c.Unlock()
}

func (c *cache) Get(key string) (interface{}, bool) {
	c.RLock()
	val, ok := c.values[key]
	c.RUnlock()
	return val, ok
}

func (c *cache) Save() error {
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

	//print(c.sid + " " + string(buf) + "\n")
	return c.cm.cli.Set(c.ctx, c.cm.getKey(c.sid), string(buf), c.expired).Err()
}

func (c *cache) Flush() error {
	c.Lock()
	c.values = make(map[string]interface{})
	c.Unlock()
	return c.Save()
}
