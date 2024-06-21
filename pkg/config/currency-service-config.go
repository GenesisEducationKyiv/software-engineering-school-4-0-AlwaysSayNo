package config

import (
	"genesis-currency-api/pkg/common/envs"
)

type CurrencyServiceConfig struct {
	ThirdPartyAPIPrivateBank  string
	ThirdPartyAPICDNJSDeliver string
}

func LoadCurrencyServiceConfig() CurrencyServiceConfig {
	return CurrencyServiceConfig{
		ThirdPartyAPIPrivateBank:  envs.Get("THIRD_PARTY_API_PRIVATE_BANK"),
		ThirdPartyAPICDNJSDeliver: envs.Get("THIRD_PARTY_API_CDN_JS_DELIVR"),
	}
}
