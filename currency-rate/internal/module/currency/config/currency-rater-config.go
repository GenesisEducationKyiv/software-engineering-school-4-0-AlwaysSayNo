package config

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-AlwaysSayNo/pkg/envs"
)

type CurrencyRaterConfig struct {
	ThirdPartyAPIPrivateBank  string
	ThirdPartyAPIBankGovUa    string
	ThirdPartyAPICDNJSDeliver string
}

func LoadCurrencyServiceConfig() CurrencyRaterConfig {
	return CurrencyRaterConfig{
		ThirdPartyAPIPrivateBank:  envs.Get("CURRENCY_SERVICE_THIRD_PARTY_API_PRIVATE_BANK"),
		ThirdPartyAPIBankGovUa:    envs.Get("CURRENCY_SERVICE_THIRD_PARTY_API_BANK_GOV_UA"),
		ThirdPartyAPICDNJSDeliver: envs.Get("CURRENCY_SERVICE_THIRD_PARTY_API_CDN_JS_DELIVR"),
	}
}
