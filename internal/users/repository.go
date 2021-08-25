package users

import (
	"basicLoginRest/internal/models"
	"context"
)

// TODO create list method with query to specify pagesize and page number
// Add the same functionality to Find
type Repository interface {
	// Find - will write found users to dest up to len(dest).
	// Returns number of written users and error.
	Find(ctx context.Context, cond models.FindUserRequest, dest []models.User) (int, error)
}
