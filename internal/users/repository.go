package users

import (
	"basicLoginRest/internal/models"
	"context"
)

type Repository interface {
	Create(ctx context.Context, u *models.User) (*models.User, error)
	Delete(ctx context.Context, userID int) error
	Update(ctx context.Context, u *models.User) (*models.User, error)

	// Find - will write found users to dest up to len(dest).
	// Returns number of written users and error.
	Find(ctx context.Context, cond models.FindUserRequest, dest []models.User) (int, error)

	GetByID(ctx context.Context, userID int) (*models.User, error)

	GetByUsername(ctx context.Context, username string, password []byte) (*models.User, error)

	GetByEmail(ctx context.Context, username string, password []byte) (*models.User, error)
}
