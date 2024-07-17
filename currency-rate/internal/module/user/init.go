package user

import (
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/api/handler"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/dto"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/model"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/repository"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/service"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/shared/dto/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user model.User) (*model.User, error)
	GetAll() (*[]model.User, error)
	ExistsByEmail(email string) bool
}

type Service interface {
	Save(saveRequestDTO dto.SaveRequestDTO) (*user.ResponseDTO, error)
	GetAll() ([]user.ResponseDTO, error)
}

type Handler interface {
	Add(ctx *gin.Context)
	FindAll(ctx *gin.Context)
}

type Module struct {
	Repository Repository
	Service    Service
	Handler    Handler
}

func Init(db *gorm.DB) *Module {
	newRepository := repository.NewRepository(db)
	userService := service.NewService(newRepository)
	userHandler := handler.NewHandler(userService)

	return &Module{
		Repository: newRepository,
		Service:    userService,
		Handler:    userHandler,
	}
}
