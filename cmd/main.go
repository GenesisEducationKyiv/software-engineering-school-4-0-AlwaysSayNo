package main

import (
	"genesis-currency-api/internal/job"
	"genesis-currency-api/internal/middleware"
	service2 "genesis-currency-api/internal/service"
	"genesis-currency-api/pkg/common/db"
	"genesis-currency-api/pkg/controller"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	dbUrl := db.GetUrl()
	d := db.Init(dbUrl)
	db.RunMigrations(dbUrl)

	r := gin.Default()
	r.Use(middleware.ErrorHandler())

	currencyService := service2.NewCurrencyService()
	controller.CurrencyRegisterRoutes(r, currencyService)

	userService := service2.NewUserService(d)
	emailService := service2.NewEmailService(userService, currencyService)
	controller.UserRegisterRoutes(r, userService, emailService)

	job.UpdateCurrency(currencyService)

	port := viper.Get("APP_PORT").(string)
	r.Run(port)
}
