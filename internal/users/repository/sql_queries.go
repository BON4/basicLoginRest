package repository

import (
	"basicLoginRest/internal/models"
	sqlr "github.com/Masterminds/squirrel"
)

var (
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
