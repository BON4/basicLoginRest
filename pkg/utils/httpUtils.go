package utils

import (
	"basicLoginRest/config"
	"basicLoginRest/internal/models"
	"basicLoginRest/pkg/httpErrors"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
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

// Configure jwt cookie
func CreateSessionCookie(cfg *config.Config, session string) *http.Cookie {
	return &http.Cookie{
		Name:  cfg.Session.Name,
		Value: session,
		Path:  "/",
		// Domain: "/",
		// Expires:    time.Now().Add(1 * time.Minute),
		RawExpires: "",
		MaxAge:     cfg.Session.Expire,
		Secure:     cfg.Cookie.Secure,
		HttpOnly:   cfg.Cookie.HTTPOnly,
		SameSite:   0,
	}
}

// Delete session
func DeleteSessionCookie(c echo.Context, sessionName string) {
	c.SetCookie(&http.Cookie{
		Name:   sessionName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
