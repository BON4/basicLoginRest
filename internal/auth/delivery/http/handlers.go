package http

import (
	"basicLoginRest/config"
	"basicLoginRest/internal/auth"
	"basicLoginRest/internal/models"
	"basicLoginRest/internal/session"
	"basicLoginRest/pkg/httpErrors"
	"basicLoginRest/pkg/logger"
	"basicLoginRest/pkg/utils"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type authHandlers struct {
	cfg *config.Config
	authUC auth.UCAuth
	sessMG session.Manager
	logger logger.Logger
	userFc models.UserFactory
}

func (a *authHandlers) Register() echo.HandlerFunc {
	type UserToRegister struct {
		Username string `json:"username"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	return func(c echo.Context) error {
		user := &UserToRegister{}
		if err := utils.ReadRequest(c, user); err != nil {
			//TODO LOG
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		valUser, err := a.userFc.NewUser(user.Username, user.Email, models.USER, user.Password)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		registerUser, err := a.authUC.Register(c.Request().Context(), &valUser)

		store, err := a.sessMG.Start(c.Request().Context(), "")
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		store.Set(registerUser.Token, registerUser.User)
		if err := store.Save(); err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, user)
	}
}

func (a *authHandlers) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(httpErrors.ErrorResponse(errors.New("not implemented")))
	}
}

func (a *authHandlers) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(httpErrors.ErrorResponse(errors.New("not implemented")))
	}
}

func (a *authHandlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(httpErrors.ErrorResponse(errors.New("not implemented")))
	}
}

func (a *authHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(httpErrors.ErrorResponse(errors.New("not implemented")))
	}
}

func (a *authHandlers) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(httpErrors.ErrorResponse(errors.New("not implemented")))
	}
}

func NewAuthHandlers(cfg *config.Config, authUC auth.UCAuth, sessMG session.Manager, uFc models.UserFactory,log logger.Logger) auth.Handlers {
	return &authHandlers{
		cfg:    cfg,
		authUC: authUC,
		sessMG: sessMG,
		logger: log,
		userFc: uFc,
	}
}