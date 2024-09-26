package config

import "github.com/GenesisEducationKyiv/software-engineering-school-4-0-AlwaysSayNo/pkg/envs"

type ServerConfig struct {
	GracefulShutdownWaitTimeSeconds int
}

func LoadServerConfigConfig() ServerConfig {
	return ServerConfig{
		GracefulShutdownWaitTimeSeconds: envs.GetInt("EMAIL_SERVICE_GRACEFUL_SHUTDOWN_WAIT_TIME_SECONDS"),
	}
}
