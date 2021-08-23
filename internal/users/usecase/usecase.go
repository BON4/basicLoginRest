package usecase

import (
	"basicLoginRest/internal/models"
	"basicLoginRest/internal/users"
	"basicLoginRest/pkg/httpErrors"
	"basicLoginRest/pkg/logger"
	"basicLoginRest/pkg/utils"
	"context"
	"github.com/pkg/errors"
)

type usersUC struct {
	usersRepo users.Repository
	logger logger.Logger
}

func (u *usersUC) LoginWithUsername(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	foundUser, err := u.usersRepo.GetByUsername(ctx, user.Username, user.Password)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.Wrap(err, "usersUC.LoginWithUsername.GetByUsername"))
	}

	token, err := utils.GenerateToken(foundUser)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "usersUC.LoginWithUsername.GenerateToken"))
	}

	return &models.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil
}

func (u *usersUC) LoginWithEmail(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	foundUser, err := u.usersRepo.GetByEmail(ctx, user.Email, user.Password)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.Wrap(err, "usersUC.LoginWithEmail.GetByEmail"))
	}

	token, err := utils.GenerateToken(foundUser)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "usersUC.LoginWithEmail.GenerateToken"))
	}

	return &models.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil
}

func (u *usersUC) Find(ctx context.Context, cond models.FindUserRequest, dest []models.User) (int, error) {
	if err := utils.ValidatePermission(ctx, models.VIEW); err != nil {
		return 0, err
	}

	if cond.PageSettings == nil {
		return 0, errors.Wrap(errors.New("page_settings condition is empty"), "usersUC.Find.EmptyCond")
	}
	//TODO could user send the empty request? It will be just select * from table;
	//if cond.Username == nil && cond.ID == nil && cond.Email == nil {
	//	return 0, errors.Wrap(errors.New("condition is empty"), "usersUC.Find.EmptyCond")
	//}

	if len(dest) == 0 {
		return 0, nil
	}

	return u.usersRepo.Find(ctx, cond, dest)
}

func (u *usersUC) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if err := utils.ValidatePermission(ctx, models.CREATE); err != nil {
		return nil, err
	}

	return u.usersRepo.Create(ctx, user)
}

func (u *usersUC) Update(ctx context.Context, user *models.User) (*models.User, error) {
	if err := utils.ValidatePermission(ctx, models.UPDATE); err != nil {
		return nil, err
	}

	_, err := u.usersRepo.GetByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return u.usersRepo.Update(ctx, user)
}

func (u *usersUC) Delete(ctx context.Context, userID int) error {
	if err := utils.ValidatePermission(ctx, models.DELETE); err != nil {
		return err
	}

	_, err := u.usersRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	return u.usersRepo.Delete(ctx, userID)
}

func (u *usersUC) GetByID(ctx context.Context, userID int) (*models.User, error) {
	if err := utils.ValidatePermission(ctx, models.VIEW); err != nil {
		return nil, err
	}

	return u.usersRepo.GetByID(ctx, userID)
}

func NewUserUseCase(repo users.Repository, log logger.Logger) users.UseCase {
	return &usersUC{
		usersRepo: repo,
		logger:    log,
	}
}