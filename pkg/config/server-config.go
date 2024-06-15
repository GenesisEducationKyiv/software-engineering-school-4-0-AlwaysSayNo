package config

import "genesis-currency-api/pkg/common/envs"

type ServerConfig struct {
	ApplicationPort string
}

func LoadServerConfigConfig() ServerConfig {
	return ServerConfig{
		ApplicationPort: envs.Get("APP_PORT"),
	}
}
