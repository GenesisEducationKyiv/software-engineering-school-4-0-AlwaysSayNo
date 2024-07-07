package main

import (
	ratecdnjsdelivr "genesis-currency-api/internal/module/currency/api/external/rater/cdnjsdelivr"
	rategovua "genesis-currency-api/internal/module/currency/api/external/rater/gov_ua"
	rateprivate "genesis-currency-api/internal/module/currency/api/external/rater/private"
	handcurrency "genesis-currency-api/internal/module/currency/api/handler"
	servcurrency "genesis-currency-api/internal/module/currency/service"
	handemail "genesis-currency-api/internal/module/email/api/handler"
	servemail "genesis-currency-api/internal/module/email/service"
	handuser "genesis-currency-api/internal/module/user/api/handler"
	repouser "genesis-currency-api/internal/module/user/repository"
	servuser "genesis-currency-api/internal/module/user/service"
	"log"

	"genesis-currency-api/pkg/config"

	"genesis-currency-api/internal/job"
	"genesis-currency-api/internal/middleware"
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
	currencyService := servcurrency.NewService(currencyProvider)
	userService := servuser.NewService(userRepository)
	emailService := servemail.NewService(userService, currencyService, config.LoadEmailServiceConfig())

	// JOBS
	job.StartAllJobs(currencyService, emailService)

	// HANDLERS
	currencyHandler := handcurrency.NewHandler(currencyService)
	handcurrency.RegisterRoutes(r, *currencyHandler)

	userHandler := handuser.NewHandler(userService)
	handuser.RegisterRoutes(r, *userHandler)

	emailHandler := handemail.NewHandler(emailService)
	handemail.RegisterRoutes(r, *emailHandler)

	// START SERVER
	cnf := config.LoadServerConfigConfig()
	if err := r.Run(cnf.ApplicationPort); err != nil {
		log.Fatal("while server bootstrapping: ", err)
	}
}

func getCurrencyProviderChain() servcurrency.Provider {
	// GET PROVIDERS
	privateClient, err := rateprivate.NewClient(config.LoadCurrencyServiceConfig())
	if err != nil {
		log.Fatal("creating Private Bank currency provider")
	}
	govUaClient, err := rategovua.NewClient(config.LoadCurrencyServiceConfig())
	if err != nil {
		log.Fatal("creating Bank Gov Ua currency provider")
	}
	jsDelivrClient, err := ratecdnjsdelivr.NewClient(config.LoadCurrencyServiceConfig())
	if err != nil {
		log.Fatal("creating JS Deliver currency provider")
	}

	// SET PROVIDERS CHAIN
	govUaClient.SetNext(jsDelivrClient)
	privateClient.SetNext(govUaClient)

	return privateClient
}
