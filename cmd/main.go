package main

import (
	"genesis-currency-api/pkg/common/db"
	"genesis-currency-api/pkg/controller"
	"genesis-currency-api/pkg/middleware"
	"genesis-currency-api/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	dbUrl := db.GetUrl()

	d := db.Init(dbUrl)

	r := gin.Default()
	r.Use(middleware.ErrorHandler())

	emailService := service.New(d)
	controller.RegisterRoutes(r, emailService)

	port := viper.Get("APP_PORT").(string)
	r.Run(port)
}
