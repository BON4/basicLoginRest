package repository

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
	store      session.UCSession
}

type Option func(*options)

func SetStore(store session.UCSession) Option {
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

type Manager struct {
	opts *options
}

func NewManager(opt ...Option) *Manager {
	opts := defaultOption

	for _, o := range opt {
		o(&opts)
	}
	return &Manager{
		opts: &opts,
	}
}

func (m *Manager) GetCookieName() string {
	return m.opts.cookieName
}

func (m *Manager) Start(ctx context.Context, sid string) (session.Session, error) {
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

func (m *Manager) Refresh(ctx context.Context, oldSid string) (session.Session, error) {
	if oldSid == "" {
		oldSid = m.opts.sessionID(ctx)
	}
	newSid := m.opts.sessionID(ctx)
	return m.opts.store.Refresh(ctx, oldSid, newSid, m.opts.expired)
}
