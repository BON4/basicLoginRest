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
	"strconv"
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

		store.Set("user", registerUser.User)
		if err := store.Save(); err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		//TODO
		//Set cookie
		c.SetCookie(utils.CreateSessionCookie(a.cfg, store.SessionID()))

		return c.JSON(http.StatusCreated, registerUser)
	}
}

func (a *authHandlers) Login() echo.HandlerFunc {
	type UserToLogin struct {
		Username string `json:"username" validate:"omitempty"`
		Email string `json:"email" validate:"omitempty"`
		Password string `json:"password"`
	}
	return func(c echo.Context) error {
		user := &UserToLogin{}
		if err := utils.ReadRequest(c, user); err != nil {
			//TODO LOG
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		var createdUser = &models.UserWithToken{}

		//If client provide zero or both credentials
		if (len(user.Username) == 0 && len(user.Email) == 0) || (len(user.Username) != 0 && len(user.Email) != 0) {
			return c.JSON(http.StatusBadRequest, httpErrors.BadQueryParams)

			//Email provided
		} else if len(user.Email) != 0 {

			valUser, err := a.userFc.NewUserWithEmail(user.Email, models.USER, user.Password)
			if err != nil {
				return c.JSON(httpErrors.ErrorResponse(err))
			}
			createdUser, err = a.authUC.LoginWithEmail(c.Request().Context(), valUser.Email, valUser.Password)
			if err != nil {
				//TODO LOG THIS
				return c.JSON(httpErrors.ErrorResponse(err))
			}

			//Username provided
		} else if len(user.Username) != 0 {

			valUser, err := a.userFc.NewUserWithUsername(user.Username, models.USER, user.Password)
			if err != nil {
				return c.JSON(httpErrors.ErrorResponse(err))
			}

			createdUser, err = a.authUC.LoginWithUsername(c.Request().Context(), valUser.Username, valUser.Password)
			if err != nil {
				//TODO LOG THIS
				return c.JSON(httpErrors.ErrorResponse(err))
			}
		}

		store, err := a.sessMG.Start(c.Request().Context(), "")
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		store.Set("user", createdUser.User)
		if err := store.Save(); err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		//Set cookie
		c.SetCookie(utils.CreateSessionCookie(a.cfg, store.SessionID()))

		return c.JSON(http.StatusCreated, createdUser)
	}
}

func (a *authHandlers) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(a.cfg.Session.Name)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(err))
			}
			return c.JSON(http.StatusInternalServerError, httpErrors.NewInternalServerError(err))
		}

		err = a.sessMG.Destroy(c.Request().Context(), cookie.Value)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		utils.DeleteSessionCookie(c, a.cfg.Session.Name)

		return c.NoContent(http.StatusOK)
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

		uid, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		user, err := a.authUC.GetByID(c.Request().Context(), uint(uid))
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, user)
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