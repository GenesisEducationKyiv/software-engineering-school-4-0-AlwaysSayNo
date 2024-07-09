package config

import (
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/pkg/common/envs"
)

type DatabaseConfig struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

func LoadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		DBUser:     envs.Get("DB_USER"),
		DBPassword: envs.Get("DB_PASSWORD"),
		DBHost:     envs.Get("DB_HOST"),
		DBPort:     envs.Get("DB_PORT"),
		DBName:     envs.Get("DB_NAME"),
	}
}
