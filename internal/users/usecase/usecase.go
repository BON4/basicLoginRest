package usecase

import (
	"basicLoginRest/internal/models"
	"basicLoginRest/internal/users"
	"basicLoginRest/pkg/logger"
	"context"
	"github.com/pkg/errors"
)

type usersUC struct {
	usersRepo users.Repository
	logger logger.Logger
}

func (u *usersUC) Find(ctx context.Context, cond models.FindUserRequest, dest []models.User) (int, error) {
	if cond.Username == nil && cond.ID == nil && cond.Email == nil {
		return 0, errors.Wrap(errors.New("condition is empty"), "usersUC.Find.EmptyCond")
	}

	if len(dest) == 0 {
		return 0, nil
	}

	return u.usersRepo.Find(ctx, cond, dest)
}

func (u *usersUC) Create(ctx context.Context, user *models.User) (*models.User, error) {
	return u.usersRepo.Create(ctx, user)
}

func (u *usersUC) Update(ctx context.Context, user *models.User) (*models.User, error) {
	_, err := u.usersRepo.GetByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return u.usersRepo.Update(ctx, user)
}

func (u *usersUC) Delete(ctx context.Context, userID int) error {
	_, err := u.usersRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	return u.usersRepo.Delete(ctx, userID)
}

func (u *usersUC) GetByID(ctx context.Context, userID int) (*models.User, error) {
	return u.usersRepo.GetByID(ctx, userID)
}

func NewUserUsecase(repo users.Repository, log logger.Logger) users.UseCase {
	return &usersUC{
		usersRepo: repo,
		logger:    log,
	}
}