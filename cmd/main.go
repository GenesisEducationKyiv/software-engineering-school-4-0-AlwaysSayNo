package main

import (
	"genesis-currency-api/internal/job"
	"genesis-currency-api/internal/middleware"
	"genesis-currency-api/internal/service"
	"genesis-currency-api/pkg/common/db"
	"genesis-currency-api/pkg/common/envs"
	"genesis-currency-api/pkg/controller"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
)

func main() {
	envs.Init()

	dbUrl := db.GetDatabaseUrl()
	d := db.Init(dbUrl)
	db.RunMigrations(dbUrl)

	r := gin.Default()
	r.Use(middleware.ErrorHandler())

	currencyService := service.NewCurrencyService()
	userService := service.NewUserService(d)
	emailService := service.NewEmailService(userService, currencyService)

	job.StartAllJobs(currencyService, emailService)
	controller.RegisterAllRoutes(r, currencyService, userService, emailService)

	port := viper.Get("APP_PORT").(string)
	err := r.Run(port)
	if err != nil {
		log.Fatal("Error happened while server bootstrapping: ", err)
	}
}
