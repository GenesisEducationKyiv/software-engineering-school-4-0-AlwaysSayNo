package user

import (
	"genesis-currency-api/internal/module/user/api/handler"
	"genesis-currency-api/internal/module/user/dto"
	"genesis-currency-api/internal/module/user/model"
	"genesis-currency-api/internal/module/user/repository"
	"genesis-currency-api/internal/module/user/service"
	"genesis-currency-api/internal/shared/dto/user"
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
