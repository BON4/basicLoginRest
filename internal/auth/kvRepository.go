package auth

import (
	"basicLoginRest/internal/models"
	"context"
)

type KVRepository interface {
	GetUserByKey(ctx context.Context, key string) (*models.User, error)
	SetUser(ctx context.Context, key string, t int, user *models.User) error
	DeleteUser(ctx context.Context, key string) error
}