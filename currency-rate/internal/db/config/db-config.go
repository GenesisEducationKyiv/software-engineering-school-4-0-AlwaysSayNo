package config

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-AlwaysSayNo/pkg/envs"
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
		DBUser:     envs.Get("CURRENCY_SERVICE_DB_USER"),
		DBPassword: envs.Get("CURRENCY_SERVICE_DB_PASSWORD"),
		DBHost:     envs.Get("CURRENCY_SERVICE_DB_HOST"),
		DBPort:     envs.Get("CURRENCY_SERVICE_DB_PORT"),
		DBName:     envs.Get("CURRENCY_SERVICE_DB_NAME"),
	}
}
