package config

import (
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/pkg/envs"
)

type ServerConfig struct {
	ApplicationPort                 string
	GracefulShutdownWaitTimeSeconds int
}

func LoadServerConfigConfig() ServerConfig {
	return ServerConfig{
		ApplicationPort:                 envs.Get("CURRENCY_SERVICE_APP_PORT"),
		GracefulShutdownWaitTimeSeconds: envs.GetInt("CURRENCY_SERVICE_GRACEFUL_SHUTDOWN_WAIT_TIME_SECONDS"),
	}
}
