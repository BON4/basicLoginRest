package repository

import (
	"basicLoginRest/internal/models"
	sqlr "github.com/Masterminds/squirrel"
)

var (
	pgGetUserByIDSqlx = func(tableName string) string {
		return `SELECT * FROM ` + tableName + ` where id = $1 LIMIT 1`
	}

	pgGetByUsernameSqlx = func(tableName string) string {
		return `SELECT * FROM ` + tableName + ` where username = $1 and password = $2 LIMIT 1`
	}

	pgGetByEmailSqlx = func(tableName string) string {
		return `SELECT * FROM ` + tableName + ` where email = $1 and password = $2 LIMIT 1`
	}

	pgCreateUserSqlx = func(tableName string) string {
		return `INSERT INTO ` + tableName + ` (username, email, role,password) values ($1, $2, $3, $4) returning *`
	}

	pgDeleteUserSqlx = func(tableName string) string {
		return `DELETE FROM ` + tableName + ` where id = $1`
	}

	pgUpdateUserSqlx = func(tableName string) string {
		return `UPDATE ` + tableName + ` set username = $1, email = $2, role = $3, password = $4 where id = $5 returning *`
	}

	pgFindUserSquirrel = func(tableName string, cond models.FindUserRequest) (string, []interface{}, error) {
		t := sqlr.Select("*").From(tableName)

		if cond.Username != nil {
			if cond.Username.Like != "" {
				t = t.Where(sqlr.Like{"username": cond.Username.Like})

			} else if cond.Username.Eq != "" {
				t = t.Where(sqlr.Eq{"username": cond.Username.Eq})
			}
		}

		if cond.Email != nil {
			if cond.Email.Like != "" {
				t = t.Where(sqlr.Like{"email": cond.Email.Like})
			} else if cond.Email.Eq != "" {
				t = t.Where(sqlr.Eq{"email": cond.Email.Eq})
			}
		}

		if cond.ID != nil {
			t = t.Where(sqlr.Eq{"id": cond.ID.Eq})
		}

		if cond.Role != nil {
			t = t.Where(sqlr.Eq{"role": cond.Role.Eq})
		}

		if cond.PageSettings != nil {
			t = t.Offset(uint64(cond.PageSettings.PageNumber))
			t = t.Limit(uint64(cond.PageSettings.PageSize))
		}

		return t.PlaceholderFormat(sqlr.Dollar).ToSql()
	}
)
