package users

import (
	"basicLoginRest/internal/models"
	"context"
)

type UseCase interface {
	LoginWithUsername(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	LoginWithEmail(ctx context.Context, user *models.User) (*models.UserWithToken, error)

	Create(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID int) error

	Find(ctx context.Context, cond models.FindUserRequest, dest []models.User) (int, error)
	GetByID(ctx context.Context, userID int) (*models.User, error)
}
