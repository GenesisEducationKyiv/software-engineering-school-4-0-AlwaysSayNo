package config

import "github.com/AlwaysSayNo/genesis-currency-api/email-service/pkg/envs"

type ServerConfig struct {
	GracefulShutdownWaitTimeSeconds int
}

func LoadServerConfigConfig() ServerConfig {
	return ServerConfig{
		GracefulShutdownWaitTimeSeconds: envs.GetInt("GRACEFUL_SHUTDOWN_WAIT_TIME_SECONDS"),
	}
}
