package http

import (
	"basicLoginRest/config"
	mock_auth "basicLoginRest/internal/auth/mock"
	"basicLoginRest/internal/models"
	mock_session "basicLoginRest/internal/session/mock"
	logger "basicLoginRest/pkg/logger"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var userFactory models.UserFactory

func AnyToBytesBuffer(i interface{}) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(i)
	if err != nil {
		return buf, err
	}
	return buf, nil
}

func TestMain(m *testing.M) {
	var err error
	userFactory, err = models.NewUserFactory(models.FactoryConfig{
		MinPasswordLen: 4,
		MinUsernameLen: 4,
		ValidateEmail:  nil,
		ParsePassword: func(password string) []byte {
			//sum := sha256.Sum256([]byte(password))
			//return sum[:]
			return sha256.New().Sum([]byte(password))
		},
	})
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestAuthHandlers_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUC := mock_auth.NewMockUCAuth(ctrl)
	mockSessMG := mock_session.NewMockManager(ctrl)

	cfg := &config.Config{
		Session: config.Session{
			Expire: 10,
		},
		Logger: config.Logger{
			Development: true,
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	handlers := NewAuthHandlers(cfg, mockAuthUC, mockSessMG, userFactory ,apiLogger)

	type UserToRegister struct {
		Username string `json:"username"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	regUser := UserToRegister{
		Username: "vlad",
		Email:    "vlad@email.com",
		Password: "1234",
	}

	buf, err := AnyToBytesBuffer(regUser)
	require.NoError(t, err)
	require.NotNil(t, buf)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(buf.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	handlerFunc := handlers.Register()

	valUser, err := userFactory.NewUser(regUser.Username, regUser.Email, models.USER, regUser.Password)
	require.NoError(t, err)

	tokenUser := &models.UserWithToken{
		User:  &valUser,
		Token: "test-token",
	}

	store := mock_session.NewMockStore(ctrl)

	mockAuthUC.EXPECT().Register(context.Background(), gomock.Eq(&valUser)).Return(tokenUser, nil)
	mockSessMG.EXPECT().Start(context.Background(), gomock.Eq("")).Return(store, nil)
	store.EXPECT().Set(gomock.Eq(tokenUser.Token), gomock.Eq(tokenUser.User))
	store.EXPECT().Save().Return(nil)
	store.EXPECT().SessionID().Return("test_sessionID")

	err = handlerFunc(c)
	require.NoError(t, err)
}

func TestAuthHandlers_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUC := mock_auth.NewMockUCAuth(ctrl)
	mockSessMG := mock_session.NewMockManager(ctrl)

	cfg := &config.Config{
		Session: config.Session{
			Expire: 10,
		},
		Logger: config.Logger{
			Development: true,
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	handlers := NewAuthHandlers(cfg, mockAuthUC, mockSessMG, userFactory ,apiLogger)

	t.Run("Login with Email", func(t *testing.T) {
		type UserToLogin struct {
			Username string `json:"username"`
			Email string `json:"email"`
			Password string `json:"password"`
		}

		logUser := UserToLogin{
			Email:    "vlad@email.com",
			Password: "1234",
		}

		buf, err := AnyToBytesBuffer(logUser)
		require.NoError(t, err)
		require.NotNil(t, buf)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		handlerFunc := handlers.Login()

		valUser, err := userFactory.NewUser("registeredUserUsername", logUser.Email, models.USER, logUser.Password)
		require.NoError(t, err)

		tokenUser := &models.UserWithToken{
			User:  &valUser,
			Token: "test-token",
		}

		store := mock_session.NewMockStore(ctrl)

		mockAuthUC.EXPECT().LoginWithEmail(context.Background(), gomock.Eq(valUser.Email), gomock.Eq(valUser.Password)).Return(tokenUser, nil)
		mockSessMG.EXPECT().Start(context.Background(), gomock.Eq("")).Return(store, nil)
		store.EXPECT().Set(gomock.Eq(tokenUser.Token), gomock.Eq(tokenUser.User))
		store.EXPECT().Save().Return(nil)
		store.EXPECT().SessionID().Return("test_sessionID")

		err = handlerFunc(c)
		require.NoError(t, err)
	})

	t.Run("Login with Username", func(t *testing.T) {
		type UserToLogin struct {
			Username string `json:"username"`
			Email string `json:"email"`
			Password string `json:"password"`
		}

		logUser := UserToLogin{
			Username:    "vlad",
			Password: "1234",
		}

		buf, err := AnyToBytesBuffer(logUser)
		require.NoError(t, err)
		require.NotNil(t, buf)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(buf.String()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		handlerFunc := handlers.Login()

		valUser, err := userFactory.NewUser(logUser.Username, "registeredUser@email.com", models.USER, logUser.Password)
		require.NoError(t, err)

		tokenUser := &models.UserWithToken{
			User:  &valUser,
			Token: "test-token",
		}

		store := mock_session.NewMockStore(ctrl)

		mockAuthUC.EXPECT().LoginWithUsername(context.Background(), gomock.Eq(valUser.Username), gomock.Eq(valUser.Password)).Return(tokenUser, nil)
		mockSessMG.EXPECT().Start(context.Background(), gomock.Eq("")).Return(store, nil)
		store.EXPECT().Set(gomock.Eq(tokenUser.Token), gomock.Eq(tokenUser.User))
		store.EXPECT().Save().Return(nil)
		store.EXPECT().SessionID().Return("test_sessionID")

		err = handlerFunc(c)
		require.NoError(t, err)
	})
}

func TestAuthHandlers_Logout(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessMG := mock_session.NewMockManager(ctrl)

	cfg := &config.Config{
		Session: config.Session{
			Name: "session-id",
			Expire: 10,
		},
		Logger: config.Logger{
			Development: true,
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	handlers := NewAuthHandlers(cfg, nil, mockSessMG, userFactory ,apiLogger)
	cookieName := cfg.Session.Name
	sessionIDStoredInCookie := "test_sid"

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.AddCookie(&http.Cookie{Name: cookieName, Value: sessionIDStoredInCookie})

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	handlerFunc := handlers.Logout()

	cookie, err := req.Cookie(cookieName)
	require.NoError(t, err)
	require.NotNil(t, cookie)

	mockSessMG.EXPECT().Destroy(context.Background(), gomock.Eq(cookie.Value)).Return(nil)
	err = handlerFunc(c)
	require.NoError(t, err)
}