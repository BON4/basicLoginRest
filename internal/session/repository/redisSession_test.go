package repository

import (
	"basicLoginRest/config"
	"basicLoginRest/internal/session"
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func getConnAndRepo() (*miniredis.Miniredis, session.Repository, *config.Config) {
	var (
		repo session.Repository
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

func TestCacheManager_Check(t *testing.T) {
	_, usecase, _ := getConnAndRepo()
	s, err := usecase.Create(context.Background(), "", time.Second)
	require.NoError(t, err)
	require.NotNil(t, s)

	err = s.Save()
	require.NoError(t, err)

	ok, err := usecase.Check(context.Background(), s.SessionID())
	require.NoError(t, err)
	require.True(t, ok)
}

func TestCacheManager_Create(t *testing.T) {
	_, usecase, _ := getConnAndRepo()
	s, err := usecase.Create(context.Background(), "", time.Second)
	require.NoError(t, err)
	require.NotNil(t, s)
}

func TestCacheManager_Update(t *testing.T) {
	_, usecase, _ := getConnAndRepo()
	s, err := usecase.Create(context.Background(), "", time.Second)
	require.NoError(t, err)
	require.NotNil(t, s)

	updatedS, err := usecase.Update(context.Background(), s.SessionID(), time.Second*5)
	require.NoError(t, err)
	require.NotNil(t, updatedS)
}

func TestCacheManager_Refresh(t *testing.T) {
	_, usecase, _ := getConnAndRepo()
	s, err := usecase.Create(context.Background(), "old", time.Second)
	require.NoError(t, err)
	require.NotNil(t, s)

	newS, err := usecase.Refresh(context.Background(), s.SessionID(), "new", time.Second*2)
	require.NoError(t, err)
	require.NotNil(t, newS)
}
