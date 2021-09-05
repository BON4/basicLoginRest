package http

import (
	"basicLoginRest/internal/auth"
	"basicLoginRest/internal/middleware"
	"basicLoginRest/internal/models"
	"github.com/labstack/echo/v4"
)

func MapAuthRoutes(authGroup *echo.Group, h auth.Handlers, mdMG *middleware.Manager) {
	authGroup.POST("/register", h.Register())
	authGroup.POST("/login", h.Login())
	authGroup.POST("/logout", h.Logout())

	authGroup.Use(mdMG.AuthSessionMiddleware)
	authGroup.GET("/:user_id", h.GetByID(), mdMG.PermissionBasedMiddleware(models.VIEW))
	authGroup.PUT("/:user_id", h.Update(), mdMG.PermissionBasedMiddleware(models.UPDATE))
	authGroup.DELETE("/:user_id", h.Delete(), mdMG.PermissionBasedMiddleware(models.DELETE))
}
