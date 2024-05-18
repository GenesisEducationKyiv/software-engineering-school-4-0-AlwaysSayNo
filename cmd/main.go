package main

import (
	"genesis-currency-api/pkg/common/db"
	"genesis-currency-api/pkg/controller"
	"genesis-currency-api/pkg/job"
	"genesis-currency-api/pkg/middleware"
	"genesis-currency-api/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	dbUrl := db.GetUrl()

	db.RunMigrations(dbUrl)
	d := db.Init(dbUrl)

	r := gin.Default()
	r.Use(middleware.ErrorHandler())

	currencyService := service.NewCurrencyService()
	controller.CurrencyRegisterRoutes(r, currencyService)

	userService := service.NewUserService(d)
	emailService := service.NewEmailService(userService, currencyService)
	controller.UserRegisterRoutes(r, userService, emailService)

	job.UpdateCurrency(currencyService)

	port := viper.Get("APP_PORT").(string)
	r.Run(port)
}
