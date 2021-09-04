package http

import (
	"basicLoginRest/internal/auth"
	"basicLoginRest/internal/middleware"
	"github.com/labstack/echo/v4"
)

func MapAuthRoutes(authGroup *echo.Group, h auth.Handlers, mdMG *middleware.Manager) {
	authGroup.POST("/register", h.Register())
	authGroup.POST("/login", h.Login())
	authGroup.POST("/logout", h.Logout())
	authGroup.GET("/:user_id", h.GetByID(), mdMG.AuthSessionMiddleware)
}
