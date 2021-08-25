package users

import (
	"basicLoginRest/internal/models"
	"context"
)

type UseCase interface {
	Find(ctx context.Context, cond models.FindUserRequest, dest []models.User) (int, error)
}
