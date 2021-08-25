package config

import (
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
)

func ParseConfig(f string) (Config, error) {
	file, err := os.Open(f)
	if err != nil {

		return Config{}, err
	}
	defer file.Close()

	opts := Config{}
	yd := yaml.NewDecoder(file)
	err = yd.Decode(&opts)

	if err != nil {
		return Config{}, err
	}
	return opts, nil
}

func ParsePostgres(f string) (PostgresConfig, error) {
	file, err := os.Open(f)
	if err != nil {

		return PostgresConfig{}, err
	}
	defer file.Close()

	opts := struct {
		Postgres PostgresConfig `yaml:"postgres"`
	}{}
	yd := yaml.NewDecoder(file)
	err = yd.Decode(&opts)

	if err != nil {
		return PostgresConfig{}, err
	}
	return opts.Postgres, nil
}

// TODO maby create a separate config for poool
func ParsePostgresPoolFromConfig(conf string) *pgxpool.Config {
	connConf := ParsePostgresConnFromConfig(conf)
	pgconf, err := pgxpool.ParseConfig("")
	if err != nil {
		panic(err)
	}
	pgconf.ConnConfig = connConf
	pgconf.MaxConns = 4
	return pgconf
}

func ParsePostgresConnFromConfig(conf string) *pgx.ConnConfig {
	postConf, err := ParsePostgres(conf)
	if err != nil {
		panic(err)
	}
	pgconf, err := pgx.ParseConfig("")
	if err != nil {
		panic(err)
	}

	values := reflect.ValueOf(postConf)
	names := reflect.Indirect(reflect.ValueOf(postConf)).Type()

	s := reflect.ValueOf(pgconf).Elem()
	for i := 0; i < values.NumField(); i++ {
		name := names.Field(i).Name
		if !(s.FieldByName(name).Type() == values.Field(i).Type()) {
			panic("Missed matched types in config struct, please specify same names as in redis.Options")
		}
		if !s.FieldByName(name).CanSet() {
			//TODO log or throw error
			//panic(fmt.Sprintf("Cant set field: %s", name))
		} else {
			s.FieldByName(name).Set(values.FieldByName(name))
		}
	}
	return pgconf
}
