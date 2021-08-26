package repository

import (
	"basicLoginRest/config"
	"basicLoginRest/internal/session"
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func getConnAndRepo() (*miniredis.Miniredis, session.UCSession, *config.Config) {
	var (
		repo session.UCSession
		opts *config.Config
	)

	mr, err := miniredis.Run()
	if err != nil {
		panic("an error '%s' was not expected when opening a stub database connection")
	}

	opts = &config.Config{Redis: config.Redis{
		Addr:            mr.Addr(),
		Database:        1,
		Password:        "123",
		MaxRetries:      3,
		MaxRetryBackoff: time.Second*5,
	}}

	mr.RequireAuth(opts.Redis.Password)

	repo = NewRedisCache(opts)

	return mr, repo, opts
}

func TestMain(m *testing.M) {
	code := m.Run()

	os.Exit(code)
}

func TestCache_Save(t *testing.T) {
	_, usecase, _ := getConnAndRepo()
	manager := NewManager(SetCookieName("go_session"), SetStore(usecase))
	s, err := manager.Start(context.Background(), "")
	require.NoError(t, err)

	key, val := "1", "2"

	s.Set(key, val)

	err = s.Save()
	require.NoError(t, err)

	savedSid := s.SessionID()

	s, err = manager.Start(context.Background(), savedSid)
	require.NoError(t, err)

	getVal, ok := s.Get(key)
	require.True(t, ok)
	require.Equal(t, getVal, val)
}
