package usecase

import (
	"basicLoginRest/internal/models"
	mock_users "basicLoginRest/internal/users/mock"
	logger "basicLoginRest/pkg/logger"
	"basicLoginRest/pkg/utils"
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

func TestUsersUC_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_users.NewMockRepository(ctrl)
	loggerApi := logger.NewApiLogger()
	userUC := NewUserUseCase(mockUserRepo, loggerApi)

	admin, _ := userFactory.NewUser("abcd", "abcd@email.com", models.ADMIN,"1234")

	userToCreate, _ := userFactory.NewUser("bcdas", "bcdas@email.com", models.VIEWER,"1234")

	ctx := context.WithValue(context.Background(), utils.UserCtxKey{}, &admin)

	mockUserRepo.EXPECT().Create(ctx, gomock.Eq(&userToCreate)).Return(&userToCreate, nil)

	createdUser, err := userUC.Create(ctx, &userToCreate)
	require.NoError(t, err)
	require.Equal(t, *createdUser, userToCreate)
}

func TestUsersUC_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_users.NewMockRepository(ctrl)
	loggerApi := logger.NewApiLogger()
	userUC := NewUserUseCase(mockUserRepo, loggerApi)

	admin, _ := userFactory.NewUser("abcd", "abcd@email.com", models.ADMIN,"1234")

	userToDelete, _ := userFactory.NewUser("bcdas", "bcdas@email.com", models.VIEWER,"1234")

	ctx := context.WithValue(context.Background(), utils.UserCtxKey{}, &admin)

	mockUserRepo.EXPECT().GetByID(ctx, gomock.Eq(userToDelete.ID)).Return(&userToDelete, nil)
	mockUserRepo.EXPECT().Delete(ctx, gomock.Eq(userToDelete.ID)).Return(nil)

	err := userUC.Delete(ctx, userToDelete.ID)
	require.NoError(t, err)
}

func TestUsersUC_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_users.NewMockRepository(ctrl)
	loggerApi := logger.NewApiLogger()
	userUC := NewUserUseCase(mockUserRepo, loggerApi)

	admin, _ := userFactory.NewUser("abcd", "abcd@email.com", models.ADMIN,"1234")

	userToUpdate, _ := userFactory.NewUser("bcdas", "bcdas@email.com", models.VIEWER,"1234")

	ctx := context.WithValue(context.Background(), utils.UserCtxKey{}, &admin)

	mockUserRepo.EXPECT().GetByID(ctx, gomock.Eq(userToUpdate.ID)).Return(&userToUpdate, nil)
	mockUserRepo.EXPECT().Update(ctx, gomock.Eq(&userToUpdate)).Return(&userToUpdate, nil)

	updatedUser, err := userUC.Update(ctx, &userToUpdate)
	require.NoError(t, err)
	require.Equal(t, updatedUser, &userToUpdate)
}

func TestUsersUC_Find(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_users.NewMockRepository(ctrl)
	loggerApi := logger.NewApiLogger()
	userUC := NewUserUseCase(mockUserRepo, loggerApi)

	admin, _ := userFactory.NewUser("abcd", "abcd@email.com", models.ADMIN,"1234")
	s := &struct {
		Like string `json:"LIKE"`
		Eq   string `json:"EQ"`
	}{Like: "cd", Eq: ""}

	p := &struct {
		PageSize uint `json:"page_size"`
		PageNumber uint `json:"page_number"`
	}{PageSize: 10, PageNumber: 1}

	findReq := models.FindUserRequest{Username: s, PageSettings: p}
	foundUsers := make([]models.User, 1)

	ctx := context.WithValue(context.Background(), utils.UserCtxKey{}, &admin)
	mockUserRepo.EXPECT().Find(ctx, gomock.Eq(findReq), foundUsers).Return(1, nil)

	n, err := userUC.Find(ctx, findReq, foundUsers)
	require.NoError(t, err)
	require.Equal(t, 1, n)
}