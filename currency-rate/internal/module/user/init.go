package user

import (
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/api/handler"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/decorator"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/repository"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler interface {
	Add(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	ChangeSubscriptionStatus(ctx *gin.Context)
}

type Module struct {
	Repository          repository.Repository
	Service             service.Service
	Handler             handler.Handler
	RepositoryDecorator decorator.Decorator
}

func Init(db *gorm.DB, producerClient decorator.ProducerClient) *Module {
	userRepository := repository.NewRepository(db)

	userService := service.NewService(userRepository)

	userServiceDecorator := decorator.NewServiceDecorator(userService, producerClient)

	userHandler := handler.NewHandler(userServiceDecorator)

	return &Module{
		*userRepository,
		*userService,
		*userHandler,
		*userServiceDecorator,
	}
}
