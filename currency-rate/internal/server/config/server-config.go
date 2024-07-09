package config

import (
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/pkg/common/envs"
)

type ServerConfig struct {
	ApplicationPort                 string
	GracefulShutdownWaitTimeSeconds int
}

func LoadServerConfigConfig() ServerConfig {
	return ServerConfig{
		ApplicationPort:                 envs.Get("APP_PORT"),
		GracefulShutdownWaitTimeSeconds: envs.GetInt("GRACEFUL_SHUTDOWN_WAIT_TIME_SECONDS"),
	}
}
