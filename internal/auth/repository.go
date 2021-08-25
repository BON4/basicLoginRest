package auth

import (
	"basicLoginRest/internal/models"
	"context"
)

type Repository interface {
	Create(ctx context.Context, u *models.User) (*models.User, error)
	Delete(ctx context.Context, userID uint) error
	Update(ctx context.Context, u *models.User) (*models.User, error)

	IsExists(ctx context.Context, username, email string) (bool, error)

	GetByID(ctx context.Context, userID uint) (*models.User, error)

	GetByUsername(ctx context.Context, username string, password []byte) (*models.User, error)

	GetByEmail(ctx context.Context, email string, password []byte) (*models.User, error)
}
