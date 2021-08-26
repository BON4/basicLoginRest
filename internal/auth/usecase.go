package auth

import (
	"basicLoginRest/internal/models"
	"context"
)

type UCAuth interface {
	LoginWithUsername(ctx context.Context, username string, password []byte) (*models.UserWithToken, error)
	LoginWithEmail(ctx context.Context, email string, password []byte) (*models.UserWithToken, error)
	Register(ctx context.Context, user *models.User) (*models.UserWithToken, error)

	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID uint) error

	GetByID(ctx context.Context, userID uint) (*models.User, error)
}