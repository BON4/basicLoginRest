package middleware

import (
	"basicLoginRest/config"
	"basicLoginRest/internal/auth"
	"basicLoginRest/internal/session"
	"basicLoginRest/pkg/logger"
)

type Manager struct {
	sessMG session.Manager
	authUC auth.UCAuth
	cfg *config.Config
	logger logger.Logger
}

func NewMiddlewareManager(sessMG session.Manager, authUC auth.UCAuth, cfg *config.Config, logger logger.Logger) *Manager {
	return &Manager{
		sessMG: sessMG,
		authUC: authUC,
		cfg:    cfg,
		logger: logger,
	}
}
