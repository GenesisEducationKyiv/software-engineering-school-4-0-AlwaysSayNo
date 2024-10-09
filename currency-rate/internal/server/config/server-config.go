package config

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-AlwaysSayNo/pkg/envs"
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
