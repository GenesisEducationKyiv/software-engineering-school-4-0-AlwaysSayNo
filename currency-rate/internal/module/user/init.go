package user

import (
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/api/handler"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/decorator"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/repository"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/service"
	"gorm.io/gorm"
)

type Module struct {
	Repository          repository.Repository
	Service             service.Service
	Handler             handler.Handler
	RepositoryDecorator decorator.Decorator
}

func Init(db *gorm.DB, producerClient decorator.ProducerClient) *Module {
	userRepository := repository.NewRepository(db)

	userService := service.NewService(userRepository)
	userHandler := handler.NewHandler(userService)

	userRepositoryDecorator := decorator.NewRepositoryDecorator(userRepository, producerClient)

	return &Module{
		*userRepository,
		*userService,
		*userHandler,
		*userRepositoryDecorator,
	}
}
