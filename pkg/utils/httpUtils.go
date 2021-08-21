package utils

import (
	"basicLoginRest/internal/models"
	"basicLoginRest/pkg/httpErrors"
	"context"
)

type UserCtxKey struct {}

func GetUserFromCtx(ctx context.Context) (*models.User, error) {
	user, ok := ctx.Value(UserCtxKey{}).(*models.User)
	if !ok {
		return nil, httpErrors.Unauthorized
	}
	return user, nil
}
