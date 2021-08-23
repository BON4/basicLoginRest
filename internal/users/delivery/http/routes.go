package http

import (
	"basicLoginRest/internal/users"
	"github.com/labstack/echo/v4"
)

func MapUsersRoutes(usersGroup *echo.Group, h users.Handlers) {
	usersGroup.POST("", h.Create())
	usersGroup.PUT("/:user_id", h.Update())
	usersGroup.DELETE("/:user_id", h.Delete())
	usersGroup.GET("/:user_id", h.GetByID())

	usersGroup.GET("/find", h.Find())
}
