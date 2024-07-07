package config

import (
	"genesis-currency-api/pkg/common/envs"
)

type CurrencyRaterConfig struct {
	ThirdPartyAPIPrivateBank  string
	ThirdPartyAPIBankGovUa    string
	ThirdPartyAPICDNJSDeliver string
}

func LoadCurrencyServiceConfig() CurrencyRaterConfig {
	return CurrencyRaterConfig{
		ThirdPartyAPIPrivateBank:  envs.Get("THIRD_PARTY_API_PRIVATE_BANK"),
		ThirdPartyAPIBankGovUa:    envs.Get("THIRD_PARTY_API_BANK_GOV_UA"),
		ThirdPartyAPICDNJSDeliver: envs.Get("THIRD_PARTY_API_CDN_JS_DELIVR"),
	}
}
