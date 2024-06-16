package controller

import (
	"genesis-currency-api/internal/service"
	"github.com/gin-gonic/gin"
)

func RegisterAllRoutes(r *gin.Engine, cs service.CurrencyService, us service.UserService, es service.EmailService) {
	RegisterCurrencyRoutes(r, cs)
	RegisterUserRoutes(r, us, es)
	RegisterUtilRoutes(r, us, es)
}
