package repository

import (
	"basicLoginRest/internal/auth"
	"basicLoginRest/internal/models"
	"basicLoginRest/pkg/dbErrors"
	"context"
	"database/sql"
	errors2 "errors"
	"github.com/jackc/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"reflect"
)

type postgresRepository struct {
	conn *sqlx.DB
	tableName string
}

func (p *postgresRepository) IsExists(ctx context.Context, username, email string) (bool, error) {
	var err error
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "pgRepository.IsExists")
		}
	}()

	exist := false

	q := pgIsUsersExists(p.tableName)
	err = p.conn.GetContext(ctx, &exist, q, username, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.Wrap(err, "pgRepository.IsExists.NotFound")
		}
		return false, errors.Wrap(err, "pgRepository.GetByCredentials")
	}

	return exist, nil
}

func (p *postgresRepository) GetByID(ctx context.Context, userID uint) (*models.User, error) {
	var err error
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "pgRepository.GetByUsername")
		}
	}()

	var pUsr models.User

	q := pgGetUserByIDSqlx(p.tableName)
	err = p.conn.GetContext(ctx, &pUsr, q, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(err, "pgRepository.GetByID.NotFound")
		}
		return nil, errors.Wrap(err, "pgRepository.GetByID")
	}

	return &pUsr, nil
}

func (p *postgresRepository) Update(ctx context.Context, u *models.User) (*models.User, error) {
	q := pgUpdateUserSqlx(p.tableName)
	rows, err := p.conn.QueryxContext(ctx, q, u.Username, u.Email, u.Role, u.Password, u.ID)
	if err != nil {
		return nil, errors.Wrap(err, "pgRepository.Update")
	}

	var pUsr models.User
	i := 0
	for rows.Next() {
		//change to simple scan for optimizing purposes
		err := rows.StructScan(&pUsr)
		if err != nil {
			return nil, errors.Wrap(err, "pgRepository.Update.StructScan")
		}
		i++
		break
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	if i == 0 {
		return nil, errors.Wrap(errors.New("User does not exists"), "pgRepository.Update.NoSuchUser")
	}
	return &pUsr, nil
}

func (p *postgresRepository) Delete(ctx context.Context, userID uint) error {
	q := pgDeleteUserSqlx(p.tableName)
	_, err := p.conn.ExecContext(ctx, q, userID)
	if err != nil {
		return errors.Wrap(err, "pgRepository.Delete")
	}
	return nil
}

func (p *postgresRepository) Create(ctx context.Context, u *models.User) (*models.User, error) {
	var perr *pgconn.PgError
	var err error
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "pgRepository.Create")
		}
	}()

	q := pgCreateUserSqlx(p.tableName)

	var createdUser models.User

	err = p.conn.QueryRowxContext(ctx, q, u.Username, u.Email, u.Role ,u.Password).
		Scan(
			&createdUser.ID,
			&createdUser.Username,
			&createdUser.Email,
			&createdUser.Role,
			&createdUser.Password,
		)

	if err != nil {
		if errors2.As(err, &perr) {
			err = dbErrors.UserAlreadyExists(perr)
		}
		return nil, err
	}

	return &createdUser, nil
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

func (p *postgresRepository) GetByUsername(ctx context.Context, username string, password []byte) (*models.User, error) {
	var err error
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "pgRepository.GetByCredentials")
		}
	}()

	var pUsr models.User

	q := pgGetByUsernameSqlx(p.tableName)
	err = p.conn.GetContext(ctx, &pUsr, q, username, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(err, "pgRepository.GetByCredentials.NotFound")
		}
		return nil, errors.Wrap(err, "pgRepository.GetByCredentials")
	}

	return &pUsr, nil
}

func (p *postgresRepository) GetByEmail(ctx context.Context, email string, password []byte) (*models.User, error) {
	var err error
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "pgRepository.GetByCredentials")
		}
	}()

	var pUsr models.User

	q := pgGetByEmailSqlx(p.tableName)
	err = p.conn.GetContext(ctx, &pUsr, q, email, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(err, "pgRepository.GetByCredentials.NotFound")
		}
		return nil, errors.Wrap(err, "pgRepository.GetByCredentials")
	}

	return &pUsr, nil
}

func NewPostgresRepository(conn *sqlx.DB, tableName string) auth.Repository {
	return &postgresRepository{conn: conn, tableName: tableName}
}
