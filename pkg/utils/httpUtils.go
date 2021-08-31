package utils

import (
	"basicLoginRest/internal/models"
	"basicLoginRest/pkg/httpErrors"
	"context"
	"github.com/labstack/echo/v4"
)

type UserCtxKey struct {}

func GetUserFromCtx(ctx context.Context) (*models.User, error) {
	user, ok := ctx.Value(UserCtxKey{}).(*models.User)
	if !ok {
		return nil, httpErrors.Unauthorized
	}
	return user, nil
}

func ReadRequest(ctx echo.Context, req interface{}) error {
	if err := ctx.Bind(req); err != nil{
		return err
	}
	return ValidateStruct(ctx.Request().Context(), req)
}