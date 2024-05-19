package controller

import (
	"genesis-currency-api/internal/service"
	"github.com/gin-gonic/gin"
)

func RegisterAllRoutes(r *gin.Engine, cs *service.CurrencyService, us *service.UserService, es *service.EmailService) {
	CurrencyRegisterRoutes(r, cs)
	UserRegisterRoutes(r, us, es)
	UtilRegisterRoutes(r, us, es)
}
