package dbErrors

import (
	"errors"
	"github.com/jackc/pgconn"
)

func UserAlreadyExists(err *pgconn.PgError) error {
	if err.SQLState() == "23505" {
		return errors.New("user with this credentials already exists")
	}
	return err
}

func UserDoesNotExists(err *pgconn.PgError) error {
	if err.SQLState() == "23505" {
		return errors.New("user with this credentials does not exists")
	}
	return err
}
