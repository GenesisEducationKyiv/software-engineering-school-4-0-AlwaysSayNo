package main

import (
	"genesis-currency-api/internal/external/api/currency/cdn_jsdelivr"
	govua "genesis-currency-api/internal/external/api/currency/gov_ua"
	"genesis-currency-api/internal/external/api/currency/private"
	repouser "genesis-currency-api/internal/repository/user"
	"log"

	"genesis-currency-api/internal/handler/currency"
	"genesis-currency-api/internal/handler/user"
	"genesis-currency-api/internal/handler/util"

	"genesis-currency-api/pkg/config"

	"genesis-currency-api/internal/job"
	"genesis-currency-api/internal/middleware"
	"genesis-currency-api/internal/service"
	"genesis-currency-api/pkg/common/db"
	"genesis-currency-api/pkg/common/envs"
	"github.com/gin-gonic/gin"
)

func main() {
	envs.Init()

	// DATABASE
	dbURL := db.GetDatabaseURL(config.LoadDatabaseConfig())
	d := db.Init(dbURL)

	// REPOSITORIES
	userRepository := repouser.NewRepository(d)

	// ENGINE
	r := gin.Default()
	r.Use(middleware.ErrorHandler())

	// CURRENCY API HANDLERS
	currencyProvider := getCurrencyProviderChain()

	// SERVICES
	currencyService := service.NewCurrencyService(currencyProvider)
	userService := service.NewUserService(userRepository)
	emailService := service.NewEmailService(userService, currencyService, config.LoadEmailServiceConfig())

	// JOBS
	job.StartAllJobs(currencyService, emailService)

	// HANDLERS
	currencyHandler := currency.NewHandler(currencyService)
	currency.RegisterRoutes(r, *currencyHandler)

	userHandler := user.NewHandler(userService)
	user.RegisterRoutes(r, *userHandler)

	utilHandler := util.NewHandler(userService, emailService)
	util.RegisterRoutes(r, *utilHandler)

	// START SERVER
	cnf := config.LoadServerConfigConfig()
	if err := r.Run(cnf.ApplicationPort); err != nil {
		log.Fatal("while server bootstrapping: ", err)
	}
}

func getCurrencyProviderChain() service.CurrencyProvider {
	// GET PROVIDERS
	privateClient, err := private.NewClient(config.LoadCurrencyServiceConfig())
	if err != nil {
		log.Fatal("creating Private Bank currency provider")
	}
	govUaClient, err := govua.NewClient(config.LoadCurrencyServiceConfig())
	if err != nil {
		log.Fatal("creating Bank Gov Ua currency provider")
	}
	jsDelivrClient, err := cdn_jsdelivr.NewClient(config.LoadCurrencyServiceConfig())
	if err != nil {
		log.Fatal("creating JS Deliver currency provider")
	}

	// SET PROVIDERS CHAIN
	govUaClient.SetNext(jsDelivrClient)
	privateClient.SetNext(govUaClient)

	return privateClient
}
