package config

import (
	"genesis-currency-api/pkg/common/envs"
)

type CurrencyServiceConfig struct {
	ThirdPartyAPI string
}

func LoadCurrencyServiceConfig() CurrencyServiceConfig {
	return CurrencyServiceConfig{
		ThirdPartyAPI: envs.Get("THIRD_PARTY_API"),
	}
}
