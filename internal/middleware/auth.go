package middleware

import (
	"basicLoginRest/internal/models"
	"basicLoginRest/pkg/httpErrors"
	"basicLoginRest/pkg/utils"
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (mg *Manager) AuthSessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(mg.cfg.Session.Name)
		if err != nil {
			//TODO LOGG
			if errors.Is(err, http.ErrNoCookie) {
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(err))
			}
			return c.JSON(http.StatusInternalServerError, httpErrors.NewInternalServerError(err))
		}

		store, err := mg.sessMG.Start(c.Request().Context(), cookie.Value)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		user := &models.User{}
		if userJson, ok := store.Get("user"); ok {
			user, ok = userJson.(*models.User)
			if !ok {
				return c.JSON(httpErrors.ErrorResponse(errors.New("no user fingerprint stored in this session")))
			}
		}

		c.Set("sid", cookie.Value)
		c.Set("role", user.Role)
		c.Set("user", user)

		ctx := context.WithValue(c.Request().Context(), utils.UserCtxKey{}, user)
		c.SetRequest(c.Request().WithContext(ctx))

		//TODO LOG THAT USER IS SUCCESSFULLY LOGGED IN
		return next(c)
	}
}

func (mg *Manager) PermissionBasedMiddleware(permissionNeeded models.Permission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole, ok := c.Get("role").(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.NoSuchRole))
			}

			if err := utils.ValidatePermission(userRole, permissionNeeded); err != nil {
				return c.JSON(err.Status(), err)
			}
			return next(c)
		}
	}
}