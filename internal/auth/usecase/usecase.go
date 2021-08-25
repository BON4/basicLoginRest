package usecase

import (
	"basicLoginRest/config"
	"basicLoginRest/internal/auth"
	"basicLoginRest/internal/models"
	"basicLoginRest/pkg/httpErrors"
	"basicLoginRest/pkg/logger"
	"basicLoginRest/pkg/utils"
	"context"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type usersUC struct {
	cfg  *config.Config
	usersRepo auth.Repository
	kvRepo auth.KVRepository
	logger logger.Logger
}

func (u *usersUC) LoginWithUsername(ctx context.Context, username string, password []byte) (*models.UserWithToken, error) {
	panic("implement me")
}

func (u *usersUC) LoginWithEmail(ctx context.Context, email string, password []byte) (*models.UserWithToken, error) {
	usr, err := u.usersRepo.GetByEmail(ctx, email, password)
	if err != nil {
		return nil, err
	}

	u.kvRepo.GetUserByID(ctx, strconv.FormatUint(uint64(usr.ID), 10))

	token, err := utils.GenerateJWTToken(u.cfg, usr)
	if err != nil {
		return nil, errors.Wrap(err, "usersUC.Register.GenerateJWTToken")
	}

	return &models.UserWithToken{
		User:  usr,
		Token: token,
	}, nil
}

func (u *usersUC) Register(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	exist, err := u.usersRepo.IsExists(ctx, user.Username, user.Email)
	if err != nil {
		return nil, err
	} else if exist {
		return nil, httpErrors.NewRestError(http.StatusBadRequest, httpErrors.ErrUserAlreadyExists, nil)
	}

	registeredUser, err := u.usersRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	
	token, err := utils.GenerateJWTToken(u.cfg, registeredUser)
	if err != nil {
		return nil, errors.Wrap(err, "usersUC.Register.GenerateJWTToken")
	}

	return &models.UserWithToken{
		User:  registeredUser,
		Token: token,
	}, nil
}

func (u *usersUC) Update(ctx context.Context, user *models.User) (*models.User, error) {
	panic("implement me")
}

func (u *usersUC) Delete(ctx context.Context, userID uint) error {
	panic("implement me")
}

func (u *usersUC) GetByID(ctx context.Context, userID uint) (*models.User, error) {
	panic("implement me")
}

func NewUserUseCase(cfg *config.Config, repo auth.Repository, kvRepo auth.KVRepository ,log logger.Logger) auth.UseCase {
	return &usersUC{
		cfg: cfg,
		usersRepo: repo,
		logger:    log,
		kvRepo: kvRepo,
	}
}
