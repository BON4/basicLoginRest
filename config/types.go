package config

import "time"

type Config struct {
	Server   ServerConfig `yaml:"server"`
	Postgres PostgresConfig `yaml:"postgres"`
	Logger   Logger `yaml:"logger"`
}

type PostgresConfig struct {
	Host           string `yaml:"host"`
	Port           uint16 `yaml:"port"`
	Database       string `yaml:"database"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	ConnectTimeout time.Duration `yaml:"connect-timeout"`
}

type ServerConfig struct {
	Host           string `yaml:"host"`
	AppVersion string `yaml:"app_version"`
	Mode string `yaml:"mode"`
	Port uint16 `yaml:"port"`
	CookieName string `yaml:"cookie_name"`
	JwtSecretKey string `yaml:"jwt_secret_key"`
}

type Logger struct {
	Development       bool `yaml:"development"`
	DisableCaller     bool `yaml:"disable_caller"`
	DisableStacktrace bool `yaml:"disable_stacktrace"`
	Encoding          string `yaml:"encoding"`
	Level             string `yaml:"level"`
}

type Cookie struct {
	MaxAgeSeconds int `yaml:"max_age_seconds"`

}