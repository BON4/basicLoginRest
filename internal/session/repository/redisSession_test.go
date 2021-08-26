package repository

import (
	"basicLoginRest/config"
	"basicLoginRest/internal/session"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"os"
	"testing"
	"time"
)

func getConnAndRepo(configPath string) (*miniredis.Miniredis, session.UCSession, *config.Config, error) {
	var (
		repo session.UCSession
		opts *config.Config
	)

	mr, err := miniredis.Run()
	if err != nil {
		return nil, nil, nil,fmt.Errorf("an error '%s' was not expected when opening a stub database connection", err)
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

	return mr, repo, opts, nil
}

func TestMain(m *testing.M) {
	code := m.Run()

	os.Exit(code)
}

func TestRedisConnect(t *testing.T) {

}
