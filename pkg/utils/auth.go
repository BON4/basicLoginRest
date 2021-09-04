package utils

import (
	"basicLoginRest/internal/models"
	"basicLoginRest/pkg/httpErrors"
)

func stringInSlice(a models.Permission, list []models.Permission) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ValidatePermission(role string, perm models.Permission) httpErrors.RestErr {
	if allowedPrm, ok := models.CheckPermission(role); ok {
		if !stringInSlice(perm,allowedPrm) {
			//TOLOG
			return httpErrors.NewForbiddenError("this user has no permission to: " + perm)
		}
	} else {
		return httpErrors.NewBadRequestError("Role does not exists")
	}

	return nil
}

func GenerateToken(user *models.User) (string, error) {
	return "not implemented token", nil
}