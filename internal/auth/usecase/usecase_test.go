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
	"strconv"
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

func TestUsersUC_LoginWithUsername(t *testing.T) {
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

	mockAuthRepo.EXPECT().GetByUsername(context.Background(), gomock.Eq(usr.Username), gomock.Eq(usr.Password)).Return(&usr, nil)

	authUser, err := authUseCase.LoginWithUsername(context.Background(), usr.Username, usr.Password)
	require.NoError(t, err)
	require.NotNil(t, authUser)
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

func TestUsersUC_GetByID(t *testing.T) {
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
	mockRedisRepo := mock_auth.NewMockKVRepository(ctrl)
	authUseCase := NewUserUseCase(&cfg, mockAuthRepo, mockRedisRepo, apiLogger)

	usr, _ := userFactory.NewUser("abcd", "abcd@email.com", models.ADMIN, "1234")
	usr.SetID(1)

	key := strconv.FormatUint(uint64(usr.ID), 10)

	mockRedisRepo.EXPECT().GetUserByKey(context.Background(), gomock.Eq(key)).Return(nil, nil)
	mockAuthRepo.EXPECT().GetByID(context.Background(), gomock.Eq(usr.ID)).Return(&usr, nil)
	mockRedisRepo.EXPECT().SetUser(context.Background(), gomock.Eq(key), gomock.Eq(cacheDurationSec), gomock.Eq(&usr)).Return(nil)

	gotUser, err := authUseCase.GetByID(context.Background(), usr.ID)
	require.NoError(t, err)
	require.NotNil(t, gotUser)
}

func TestUsersUC_Update(t *testing.T) {
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
	mockRedisRepo := mock_auth.NewMockKVRepository(ctrl)
	authUseCase := NewUserUseCase(&cfg, mockAuthRepo, mockRedisRepo, apiLogger)

	usr, _ := userFactory.NewUser("abcd", "abcd@email.com", models.ADMIN, "1234")
	usr.SetID(1)

	key := strconv.FormatUint(uint64(usr.ID), 10)

	mockAuthRepo.EXPECT().Update(context.Background(), gomock.Eq(&usr)).Return(&usr, nil)
	mockRedisRepo.EXPECT().DeleteUser(context.Background(), gomock.Eq(key)).Return(nil)

	updatedUser, err := authUseCase.Update(context.Background(), &usr)
	require.NoError(t, err)
	require.NotNil(t, updatedUser)
}

func TestUsersUC_Delete(t *testing.T) {
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
	mockRedisRepo := mock_auth.NewMockKVRepository(ctrl)
	authUseCase := NewUserUseCase(&cfg, mockAuthRepo, mockRedisRepo, apiLogger)

	usr, _ := userFactory.NewUser("abcd", "abcd@email.com", models.ADMIN, "1234")
	usr.SetID(1)

	key := strconv.FormatUint(uint64(usr.ID), 10)

	mockAuthRepo.EXPECT().Delete(context.Background(), gomock.Eq(usr.ID)).Return(nil)
	mockRedisRepo.EXPECT().DeleteUser(context.Background(), gomock.Eq(key)).Return(nil)

	err := authUseCase.Delete(context.Background(), usr.ID)
	require.NoError(t, err)
}