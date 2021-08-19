package config

import "time"

type DatabaseConfig struct {
	PgConfig *struct {
		Host           string `yaml:"host"`
		Port           uint16 `yaml:"port"`
		Database       string `yaml:"database"`
		User           string `yaml:"user"`
		Password       string `yaml:"password"`
		ConnectTimeout time.Duration `yaml:"connect-timeout"`
	} `yaml:"postgres"`
}
