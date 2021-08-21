package users

import (
	"basicLoginRest/internal/models"
	"context"
)

// TODO create list method with query to specify pagesize and page number
// Add the same functionality to Find
type UseCase interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID int) error
	Find(ctx context.Context, cond models.FindUserRequest, dest []models.User) (int, error)
	GetByID(ctx context.Context, userID int) (*models.User, error)
}
