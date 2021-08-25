package usecase

import (
	"basicLoginRest/config"
	mock_auth "basicLoginRest/internal/auth/mock"
	"basicLoginRest/internal/models"
	"basicLoginRest/pkg/logger"
	"context"
	"crypto/sha256"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

var userFactory models.UserFactory

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

func TestUsersUC_Register(t *testing.T) {
	cfg := config.Config{
		Server:   config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger:   config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(&cfg)
	mockAuthRepo := mock_auth.NewMockRepository(ctrl)
	authUseCase := NewUserUseCase(&cfg, mockAuthRepo, nil, apiLogger)

	usr, _ := userFactory.NewUser("abcd", "abcd@email.com", models.ADMIN, "1234")

	mockAuthRepo.EXPECT().IsExists(context.Background(), gomock.Eq(usr.Username), gomock.Eq(usr.Email)).Return(false, nil)
	mockAuthRepo.EXPECT().Create(context.Background(), gomock.Eq(&usr)).Return(&usr, nil)

	registeredUser, err := authUseCase.Register(context.Background(), &usr)
	require.NoError(t, err)
	require.NotNil(t, registeredUser)
}

func TestUsersUC_LoginWithEmail(t *testing.T) {
	cfg := config.Config{
		Server:   config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger:   config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(&cfg)
	mockAuthRepo := mock_auth.NewMockRepository(ctrl)
	authUseCase := NewUserUseCase(&cfg, mockAuthRepo, nil, apiLogger)

	usr, _ := userFactory.NewUser("abcd", "abcd@email.com", models.ADMIN, "1234")
	usr.SetID(1)

	mockAuthRepo.EXPECT().GetByEmail(context.Background(), gomock.Eq(usr.Email), gomock.Eq(usr.Password)).Return(&usr, nil)

	authUser, err := authUseCase.LoginWithEmail(context.Background(), usr.Email, usr.Password)
	require.NoError(t, err)
	require.NotNil(t, authUser)
}