package utils

import (
	"basicLoginRest/internal/models"
	"context"
	"github.com/pkg/errors"
)

func stringInSlice(a models.Permission, list []models.Permission) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ValidatePermission(ctx context.Context, perm models.Permission) error {
	u, err := GetUserFromCtx(ctx)
	if err != nil {
		return err
	}

	if allowedPrm, ok := models.CheckPermission(u.Role); ok {
		if !stringInSlice(perm,allowedPrm) {
			//TOLOG
			return errors.New("Method not allowed for this role")
		}
	} else {
		return errors.New("Role does not exists")
	}

	return nil
}

func GenerateToken(user *models.User) (string, error) {
	return "not implemented token", nil
}