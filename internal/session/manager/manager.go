//TODO change options type of initializeing to simple constructor with confg.Config
package manager

import (
	"basicLoginRest/internal/session"
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

type IDHandlerFunc func(context.Context) string

func newUUID () string {
	//var buf [10]byte
	//io.ReadFull(rand.Reader, buf[:])
	//
	//return string(buf[:])
	//return fmt.Sprintf("%d", rand.Int())
	return uuid.NewString()
}

var defaultOption = options{
	sessionID: func(_ context.Context) string {
		return newUUID()
	},
	cookieName: "go_session_id",
	expired:    time.Minute,
	store:      nil,
}

type options struct {
	sessionID  IDHandlerFunc
	cookieName string
	expired    time.Duration
	store      session.Repository
}

type Option func(*options)

func SetStore(store session.Repository) Option {
	return func(o *options) {
		o.store = store
	}
}

func SetExpired(exp time.Duration) Option {
	return func(o *options) {
		o.expired = exp
	}
}

func SetCookieName(name string) Option {
	return func(o *options) {
		o.cookieName = name
	}
}

type manager struct {
	opts *options
}

func newManager(opt ...Option) session.Manager {
	opts := defaultOption

	for _, o := range opt {
		o(&opts)
	}
	return &manager{
		opts: &opts,
	}
}

func (m *manager) GetCookieName() string {
	return m.opts.cookieName
}

func (m *manager) Start(ctx context.Context, sid string) (session.Store, error) {
	//TODO Maybe create function validSid
	if sid != "" {
		ok, err := m.opts.store.Check(ctx, sid)
		if err != nil {
			return nil, err
		} else if ok {
			return m.opts.store.Update(ctx, sid, m.opts.expired)
		} else {
			return nil, errors.New("session not found")
		}
	}

	newSid := m.opts.sessionID(ctx)
	//TODO Set this sid in to cookie
	return m.opts.store.Create(ctx, newSid, m.opts.expired)
}

func (m *manager) Refresh(ctx context.Context, oldSid string) (session.Store, error) {
	if oldSid == "" {
		oldSid = m.opts.sessionID(ctx)
	}
	newSid := m.opts.sessionID(ctx)
	return m.opts.store.Refresh(ctx, oldSid, newSid, m.opts.expired)
}

func (m *manager) Destroy(ctx context.Context, sid string) error {
	if sid == "" {
		return nil
	}
	return m.opts.store.Delete(ctx, sid)
}
