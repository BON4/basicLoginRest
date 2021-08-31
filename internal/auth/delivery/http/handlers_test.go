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
	mockSessMG.EXPECT().Start(context.Background(), gomock.Eq("test-sid")).Return(store, nil)
	store.EXPECT().Set(gomock.Eq(tokenUser.Token), gomock.Eq(tokenUser.User))
	store.EXPECT().Save().Return(nil)

	err = handlerFunc(c)
	require.NoError(t, err)
}