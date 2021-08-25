package repository

import (
	"basicLoginRest/internal/models"
	"basicLoginRest/internal/users"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"reflect"
)

type postgresRepository struct {
	conn *sqlx.DB
	tableName string
}


func (p *postgresRepository) Find(ctx context.Context, cond models.FindUserRequest, dest []models.User) (int, error) {
	//TODO refactor this method with sqlx.Select for compatibility. Although allocations can rise little bit
	var err error
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "pgRepository.Find")
		}
	}()

	sqlStr, args, err := pgFindUserSquirrel(p.tableName, cond)
	if err != nil {
		return 0, errors.Wrap(err, "pgRepository.Find.FindBuilder")
	}

	rows, err := p.conn.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var usr models.User
	s := reflect.ValueOf(&usr).Elem()
	numCols := s.NumField()
	columns := make([]interface{}, numCols)
	for i := 0; i < numCols; i++ {
		field := s.Field(i)
		columns[i] = field.Addr().Interface()
	}

	i := 0
	for rows.Next() {
		if i >= len(dest) {
			return i, nil
		}

		err := rows.Scan(columns...)
		if err != nil {
			return 0, err
		}

		dest[i].Username = usr.Username
		dest[i].Email = usr.Email
		dest[i].Password = usr.Password
		dest[i].ID = usr.ID
		dest[i].Role = usr.Role

		i++
	}

	if rows.Err() != nil {
		return 0, rows.Err()
	}

	return i, nil
}

func NewPostgresRepository(conn *sqlx.DB, tableName string) users.Repository {
	return &postgresRepository{conn: conn, tableName: tableName}
}