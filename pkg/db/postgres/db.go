package postgres

import (
	"basicLoginRest/config"
	"context"
	"database/sql"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"os"
	"sync"
)

var (
	once sync.Once
	oncePool sync.Once

	DB *sqlx.DB
	Pool *pgxpool.Pool
)

type migrateFunc func(ctx context.Context, DB *sql.DB) error

//OpenSqlxViaPgx - opens db connection with pgx driver, provide migration function if you want to create migration
func OpenSqlxViaPgx(ctx context.Context, configF string, mf migrateFunc) *sqlx.DB {
	once.Do(func() {
		conf := config.ParsePostgresConnFromConfig(configF)
		//TODO create connection pool config parser
		nativeDB := stdlib.OpenDB(*conf)

		//Apply migrations
		if mf != nil {
			err := mf(ctx, nativeDB)
			if err != nil {
				panic(errors.Wrap(err, "error while applying migrations"))
			}
		}
		DB = sqlx.NewDb(nativeDB, "pgx")
	})
	return DB
}

func OpenSqlxViaSql(ctx context.Context, db *sql.DB, driverName string, mf migrateFunc) *sqlx.DB {
	once.Do(func() {
		if mf != nil {
			err := mf(ctx, db)
			if err != nil {
				panic(errors.Wrap(err, "error while applying migrations"))
			}
		}

		DB = sqlx.NewDb(db, driverName)
	})
	return DB
}

func Migrate(ctx context.Context, DB *sql.DB, migrateFile string) error {
	f, err := os.Open(migrateFile)
	if err != nil {
		return err
	}
	defer f.Close()

	fs, err := f.Stat()
	if err != nil {
		return err
	}

	b := make([]byte, fs.Size())
	n, err := f.Read(b)
	if err != nil {
		return err
	}

	_, err = DB.ExecContext(ctx, string(b[:n]))
	return err
}