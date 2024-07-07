package main

import (
	"genesis-currency-api/internal/db"
	config4 "genesis-currency-api/internal/db/config"
	currencymodule "genesis-currency-api/internal/module/currency"
	ratecdnjsdelivr "genesis-currency-api/internal/module/currency/api/external/rater/cdnjsdelivr"
	rategovua "genesis-currency-api/internal/module/currency/api/external/rater/gov_ua"
	rateprivate "genesis-currency-api/internal/module/currency/api/external/rater/private"
	currencyhand "genesis-currency-api/internal/module/currency/api/handler"
	currencyconf "genesis-currency-api/internal/module/currency/config"
	currencyserv "genesis-currency-api/internal/module/currency/service"
	emailmodule "genesis-currency-api/internal/module/email"
	emailhand "genesis-currency-api/internal/module/email/api/handler"
	emailconf "genesis-currency-api/internal/module/email/config"
	usermodule "genesis-currency-api/internal/module/user"
	userhand "genesis-currency-api/internal/module/user/api/handler"
	userconf "genesis-currency-api/internal/server/config"
	"log"

	"genesis-currency-api/internal/job"
	"genesis-currency-api/internal/middleware"
	"genesis-currency-api/pkg/common/envs"
	"github.com/gin-gonic/gin"
)

func main() {
	envs.Init()

	// DATABASE
	dbURL := db.GetDatabaseURL(config4.LoadDatabaseConfig())
	d := db.Init(dbURL)

	// CURRENCY API HANDLERS
	currencyProvider := getCurrencyProviderChain()

	// MODULES
	userModule := usermodule.Init(d)
	currencyModule := currencymodule.Init(currencyProvider)
	emailModule := emailmodule.Init(userModule.Service, currencyModule.Service, emailconf.LoadEmailServiceConfig())

	// ENGINE
	r := gin.Default()
	r.Use(middleware.ErrorHandler())

	// JOBS
	job.StartAllJobs(currencyModule.Service, emailModule.Service)

	// HANDLERS
	currencyhand.RegisterRoutes(r, currencyModule.Handler)
	userhand.RegisterRoutes(r, userModule.Handler)
	emailhand.RegisterRoutes(r, emailModule.Handler)

	// START SERVER
	cnf := userconf.LoadServerConfigConfig()
	if err := r.Run(cnf.ApplicationPort); err != nil {
		log.Fatal("while server bootstrapping: ", err)
	}
}

func getCurrencyProviderChain() currencyserv.Provider {
	// GET PROVIDERS
	privateClient, err := rateprivate.NewClient(currencyconf.LoadCurrencyServiceConfig())
	if err != nil {
		log.Fatal("creating Private Bank currency provider")
	}
	govUaClient, err := rategovua.NewClient(currencyconf.LoadCurrencyServiceConfig())
	if err != nil {
		log.Fatal("creating Bank Gov Ua currency provider")
	}
	jsDelivrClient, err := ratecdnjsdelivr.NewClient(currencyconf.LoadCurrencyServiceConfig())
	if err != nil {
		log.Fatal("creating JS Deliver currency provider")
	}

	// SET PROVIDERS CHAIN
	govUaClient.SetNext(jsDelivrClient)
	privateClient.SetNext(govUaClient)

	return privateClient
}
