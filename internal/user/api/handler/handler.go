package handler

import (
	"net/http"

	"genesis-currency-api/pkg/dto"
	"genesis-currency-api/pkg/errors"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	Save(user dto.UserSaveRequestDTO) (*dto.UserResponseDTO, error)
	GetAll() ([]dto.UserResponseDTO, error)
}

type Handler struct {
	userService UserService
}

func NewHandler(saver UserService) *Handler {
	return &Handler{
		userService: saver,
	}
}

func (h *Handler) Add(ctx *gin.Context) {
	email := ctx.PostForm("email")

	var saveDto dto.UserSaveRequestDTO
	saveDto.Email = email

	_, err := h.userService.Save(saveDto)
	if err != nil {
		errors.AttachToCtx(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, "E-mail додано")
}

func (h *Handler) FindAll(ctx *gin.Context) {
	result, err := h.userService.GetAll()
	if err != nil {
		errors.AttachToCtx(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, &result)
}

// RegisterRoutes creates an instance of Handler and registers routes for it.
func RegisterRoutes(r *gin.Engine, handler Handler) {
	routes := r.Group("/api/v1/")
	routes.POST("/subscribe", handler.Add)
	routes.GET("/users", handler.FindAll)
}
