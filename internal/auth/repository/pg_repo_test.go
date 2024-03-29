package repository

import (
	"basicLoginRest/internal/models"
	"basicLoginRest/pkg/db/postgres"
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	userFactory models.UserFactory
	skipDatabaseTest bool = false
	DB *sqlx.DB
)

const tableName = "usersTest"

func TestMain(m *testing.M) {
	var err error
	userFactory, err = models.NewUserFactory(models.FactoryConfig{
		MinPasswordLen: 4,
		MinUsernameLen: 4,
		ValidateEmail:  nil,
		ParsePassword: func(password string) []byte {
			//sum := sha256.Sum256([]byte(password))
			//return sum[:]
			return sha256.New().Sum([]byte(password))
		},
	})
	if err != nil {
		panic(err)
	}

	DB = postgres.OpenSqlxViaPgx(context.Background(), "C:\\Users\\home\\go\\src\\basicLoginRest\\config\\config.yaml", nil)
	defer DB.Close()

	err = DB.Ping()
	if err != nil {
		skipDatabaseTest = true
	}
	m.Run()
}

func TestMockGetByCredentials(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	repo := NewPostgresRepository(sqlxDB, tableName)

	t.Run("GetByUsername", func(t *testing.T) {
		u, _ := userFactory.NewUser("abcd", "abcd@gmail.com", models.USER,"1324")

		rows := sqlmock.NewRows([]string{"id", "username", "email", "role", "password"}).AddRow(u.ID, u.Username, u.Email, u.Role, u.Password)

		q := pgGetByUsernameSqlx(tableName)
		//q, args := createUserJet(&u)
		mock.ExpectQuery(q).WithArgs(u.Username, u.Password).WillReturnRows(rows)

		createdUser, err := repo.GetByUsername(context.Background(), u.Username, u.Password)

		require.NoError(t, err)
		require.NotNil(t, createdUser)
		require.Equal(t, createdUser, &u)
	})

	t.Run("GetByUsername Error", func(t *testing.T) {
		u, _ := userFactory.NewUser("abcd", "abcd@gmail.com", models.USER,"1324")

		q := pgGetByUsernameSqlx(tableName)
		//q, args := createUserJet(&u)
		mock.ExpectQuery(q).WithArgs(u.Username, u.Password).WillReturnError(sql.ErrNoRows)

		createdUser, err := repo.GetByUsername(context.Background(), u.Username, u.Password)

		require.Error(t, err)
		require.Nil(t, createdUser)
	})

	t.Run("GetByEmail", func(t *testing.T) {
		u, _ := userFactory.NewUser("abcd", "abcd@gmail.com", models.USER,"1324")

		rows := sqlmock.NewRows([]string{"id", "username", "email", "role", "password"}).AddRow(u.ID, u.Username, u.Email, u.Role, u.Password)

		q := pgGetByEmailSqlx(tableName)
		//q, args := createUserJet(&u)
		mock.ExpectQuery(q).WithArgs(u.Email, u.Password).WillReturnRows(rows)

		createdUser, err := repo.GetByEmail(context.Background(), u.Email, u.Password)

		require.NoError(t, err)
		require.NotNil(t, createdUser)
		require.Equal(t, createdUser, &u)
	})

	t.Run("GetByEmail Error", func(t *testing.T) {
		u, _ := userFactory.NewUser("abcd", "abcd@gmail.com", models.USER,"1324")

		q := pgGetByEmailSqlx(tableName)
		//q, args := createUserJet(&u)
		mock.ExpectQuery(q).WithArgs(u.Email, u.Password).WillReturnError(sql.ErrNoRows)

		createdUser, err := repo.GetByUsername(context.Background(), u.Email, u.Password)

		require.Error(t, err)
		require.Nil(t, createdUser)
	})
}

func TestMockUpdate(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	repo := NewPostgresRepository(sqlxDB, tableName)
	
	t.Run("Update", func(t *testing.T) {
		u, _ := userFactory.NewUser("abcd", "abcd@gmail.com", models.USER,"1324")

		rows := sqlmock.NewRows([]string{"id", "username", "email", "role", "password"}).AddRow(u.ID, u.Username, u.Email, u.Role, u.Password)

		q := pgUpdateUserSqlx(tableName)
		//q, args := createUserJet(&u)
		mock.ExpectQuery(q).WithArgs(u.Username, u.Email, u.Role, u.Password, u.ID).WillReturnRows(rows)

		createdUser, err := repo.Update(context.Background(), &u)

		require.NoError(t, err)
		require.NotNil(t, createdUser)
		require.Equal(t, createdUser, &u)
	})

	t.Run("Update Does Not Exists", func(t *testing.T) {
		u, _ := userFactory.NewUser("abcd", "abcd@gmail.com", models.USER,"1324")

		q := pgUpdateUserSqlx(tableName)

		//No rows will be returned if updating user does not exists
		rows := sqlmock.NewRows([]string{})
		//q, args := createUserJet(&u)
		mock.ExpectQuery(q).WithArgs(u.Username, u.Email, u.Role, u.Password, u.ID).WillReturnRows(rows)

		updatedUser, err := repo.Update(context.Background(), &u)
		require.Error(t, err)
		require.Nil(t, updatedUser)
	})
}

func TestMockCreate(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	repo := NewPostgresRepository(sqlxDB, tableName)

	t.Run("Create", func(t *testing.T) {
		u, _ := userFactory.NewUser("abcd", "abcd@gmail.com", models.USER,"1324")

		rows := sqlmock.NewRows([]string{"id", "username", "email", "role", "password"}).AddRow(u.ID, u.Username, u.Email, u.Role, u.Password)

		q := pgCreateUserSqlx(tableName)
		//q, args := createUserJet(&u)
		mock.ExpectQuery(q).WithArgs(u.Username, u.Email, u.Role, u.Password).WillReturnRows(rows)

		createdUser, err := repo.Create(context.Background(), &u)

		require.NoError(t, err)
		require.NotNil(t, createdUser)
		require.Equal(t, createdUser, &u)
	})

	t.Run("Create Err Already Exists", func(t *testing.T) {
		u, _ := userFactory.NewUser("abcd", "abcd@gmail.com", models.USER,"1324")

		expectedErr := errors.New("user with this credentials already exists")

		q := pgCreateUserSqlx(tableName)

		//q, args := createUserJet(&u)
		mock.ExpectQuery(q).WithArgs(u.Username, u.Email, u.Role, u.Password).WillReturnError(expectedErr)

		//Then attempt to create user with the same credentials
		createdUser, err := repo.Create(context.Background(), &u)

		require.Nil(t, createdUser)
		require.NotNil(t, err)
	})
}

func TestPostgresRepository_Create(t *testing.T) {
	if skipDatabaseTest {
		t.Skip("No connection to database")
		return
	}

	defer func() {
		_, err := DB.ExecContext(context.Background(), "DELETE FROM " + tableName)
		if err != nil {
			panic(err)
		}
	}()

	userrepo := NewPostgresRepository(DB, tableName)

	t.Run("Create", func(t *testing.T) {
		usr, _ := userFactory.NewUser("abcd", "abcd@gmail.com", models.USER,"1324")
		_, err := userrepo.Create(context.Background(), &usr)
		require.NoError(t, err)
	})

	t.Run("Create Error", func(t *testing.T) {
		usr, _ := userFactory.NewUser("abcd", "abcd@gmail.com", models.USER,"1324")
		_, err := userrepo.Create(context.Background(), &usr)
		require.Error(t, err)
		fmt.Println(err)
	})
}

func TestPostgresRepository_Delete(t *testing.T) {
	if skipDatabaseTest {
		t.Skip("No connection to database")
		return
	}

	defer func() {
		_, err := DB.ExecContext(context.Background(), "DELETE FROM " + tableName)
		if err != nil {
			panic(err)
		}
	}()

	userrepo := NewPostgresRepository(DB, tableName)

	t.Run("Delete", func(t *testing.T) {
		usr, _ := userFactory.NewUser("abcd", "abcd@gmail.com", models.USER,"1324")
		newUser, err := userrepo.Create(context.Background(), &usr)
		require.NoError(t, err)

		err = userrepo.Delete(context.Background(), newUser.ID)
		require.NoError(t, err)

		_, err = userrepo.GetByID(context.Background(), newUser.ID)
		require.Error(t, err)
	})
}

func TestPostgresRepository_GetByCredentials(t *testing.T) {
	if skipDatabaseTest {
		t.Skip("No connection to database")
		return
	}

	defer func() {
		_, err := DB.ExecContext(context.Background(), "DELETE FROM " + tableName)
		if err != nil {
			panic(err)
		}
	}()

	userrepo := NewPostgresRepository(DB, tableName)

	usr, _ := userFactory.NewUser("abcd", "abcd@gmail.com", models.USER,"1324")
	newUser, err := userrepo.Create(context.Background(), &usr)
	require.NoError(t, err)

	t.Run("GetByUsername", func(t *testing.T) {
		_, err = userrepo.GetByUsername(context.Background(), newUser.Username, newUser.Password)
		require.NoError(t, err)
	})

	t.Run("GetByUsername Error", func(t *testing.T) {
		_, err = userrepo.GetByUsername(context.Background(), newUser.Username + "offset", newUser.Password)
		require.Error(t, err)
	})

	t.Run("GetByEmail", func(t *testing.T) {
		_, err = userrepo.GetByEmail(context.Background(), newUser.Email, newUser.Password)
		require.NoError(t, err)
	})

	t.Run("GetByEmail Error", func(t *testing.T) {
		_, err = userrepo.GetByEmail(context.Background(), newUser.Email + "offset", newUser.Password)
		require.Error(t, err)
	})
}

func TestPostgresRepository_IsExists(t *testing.T) {
	if skipDatabaseTest {
		t.Skip("No connection to database")
		return
	}

	defer func() {
		_, err := DB.ExecContext(context.Background(), "DELETE FROM " + tableName)
		if err != nil {
			panic(err)
		}
	}()

	userrepo := NewPostgresRepository(DB, tableName)

	usr, _ := userFactory.NewUser("abcd", "abcd@gmail.com", models.USER,"1324")
	newUser, err := userrepo.Create(context.Background(), &usr)
	require.NoError(t, err)

	t.Run("IsExists", func(t *testing.T) {
		ok, err := userrepo.IsExists(context.Background(), newUser.Username, newUser.Email)
		require.NoError(t, err)
		require.True(t, ok)

		ok, err = userrepo.IsExists(context.Background(), newUser.Username + "offset", newUser.Email)
		require.NoError(t, err)
		require.True(t, ok)

		ok, err = userrepo.IsExists(context.Background(), newUser.Username, newUser.Email + "offset")
		require.NoError(t, err)
		require.True(t, ok)

		ok, err = userrepo.IsExists(context.Background(), newUser.Username + "offset", newUser.Email + "offset")
		require.NoError(t, err)
		require.False(t, ok)
	})
}

func TestPostgresRepository_Update(t *testing.T) {
	if skipDatabaseTest {
		t.Skip("No connection to database")
		return
	}

	defer func() {
		_, err := DB.ExecContext(context.Background(), "DELETE FROM " + tableName)
		if err != nil {
			panic(err)
		}
	}()

	userrepo := NewPostgresRepository(DB, tableName)

	t.Run("Update", func(t *testing.T) {
		usr, _ := userFactory.NewUser("abcd", "abcd@gmail.com", models.USER,"1324")
		newUser, err := userrepo.Create(context.Background(), &usr)
		require.NoError(t, err)

		preparedUser, err := userFactory.NewUser("efgk", "efgk@gmail.com", models.USER,"1234")
		require.NoError(t, err)
		preparedUser.SetID(newUser.ID)

		updatedUser, err := userrepo.Update(context.Background(), &preparedUser)
		require.NoError(t, err)
		require.Equal(t, updatedUser.ID, newUser.ID)
		require.Equal(t, updatedUser.Username, preparedUser.Username)
		require.Equal(t, updatedUser.Email, preparedUser.Email)
		require.Equal(t, updatedUser.Password, preparedUser.Password)
	})
}