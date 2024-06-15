package main

import (
	"log"

	"genesis-currency-api/pkg/config"

	"genesis-currency-api/internal/job"
	"genesis-currency-api/internal/middleware"
	"genesis-currency-api/internal/service"
	"genesis-currency-api/pkg/common/db"
	"genesis-currency-api/pkg/common/envs"
	"genesis-currency-api/pkg/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	envs.Init()

	dbURL := db.GetDatabaseURL(config.LoadDatabaseConfig())
	d := db.Init(dbURL)
	db.RunMigrations(d)

	r := gin.Default()
	r.Use(middleware.ErrorHandler())

	currencyService := service.NewCurrencyService(config.LoadCurrencyServiceConfig())
	userService := service.NewUserService(d)
	emailService := service.NewEmailService(userService, currencyService, config.LoadEmailServiceConfig())

	job.StartAllJobs(currencyService, emailService)
	controller.RegisterAllRoutes(r, currencyService, userService, emailService)

	cnf := config.LoadServerConfigConfig()
	if err := r.Run(cnf.ApplicationPort); err != nil {
		log.Fatal("error happened while server bootstrapping: ", err)
	}
}
