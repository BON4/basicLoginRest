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

var cacheDurationSec = 60

func (u *usersUC) LoginWithUsername(ctx context.Context, username string, password []byte) (*models.UserWithToken, error) {
	usr, err := u.usersRepo.GetByUsername(ctx, username, password)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateJWTToken(u.cfg, usr)
	if err != nil {
		return nil, errors.Wrap(err, "usersUC.LoginWithUsername.GenerateJWTToken")
	}

	return &models.UserWithToken{
		User:  usr,
		Token: token,
	}, nil
}

func (u *usersUC) LoginWithEmail(ctx context.Context, email string, password []byte) (*models.UserWithToken, error) {
	usr, err := u.usersRepo.GetByEmail(ctx, email, password)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateJWTToken(u.cfg, usr)
	if err != nil {
		return nil, errors.Wrap(err, "usersUC.LoginWithEmail.GenerateJWTToken")
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
	key := strconv.FormatUint(uint64(user.ID), 10)
	updatedUser, err := u.usersRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	if err := u.kvRepo.DeleteUser(ctx, key); err != nil {
		u.logger.Errorf("usersUC.Update.DeleteUser: %v", err)
	}

	return updatedUser, nil
}

func (u *usersUC) Delete(ctx context.Context, userID uint) error {
	key := strconv.FormatUint(uint64(userID), 10)

	if err := u.usersRepo.Delete(ctx, userID); err != nil {
		return err
	}

	if err := u.kvRepo.DeleteUser(ctx, key); err != nil {
		u.logger.Errorf("usersUC.Update.DeleteUser: %v", err)
	}

	return nil
}

func (u *usersUC) GetByID(ctx context.Context, userID uint) (*models.User, error) {
	key := strconv.FormatUint(uint64(userID), 10)
	redisUser, err := u.kvRepo.GetUserByKey(ctx,key )
	if err != nil {
		return nil, err
	}
	if redisUser != nil {
		return redisUser, nil
	}

	repoUser, err := u.usersRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := u.kvRepo.SetUser(ctx, key, cacheDurationSec, repoUser); err != nil {
		u.logger.Errorf("usersUC.GetByID.SetUser: %v", err)
	}

	return repoUser, nil
}

func NewUserUseCase(cfg *config.Config, repo auth.Repository, kvRepo auth.KVRepository ,log logger.Logger) auth.UCAuth {
	return &usersUC{
		cfg: cfg,
		usersRepo: repo,
		logger:    log,
		kvRepo: kvRepo,
	}
}
