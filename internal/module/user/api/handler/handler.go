package handler

import (
	"genesis-currency-api/internal/module/user/dto"
	"genesis-currency-api/internal/shared/dto/user"
	"net/http"

	"genesis-currency-api/pkg/errors"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	Save(user dto.SaveRequestDTO) (*user.ResponseDTO, error)
	GetAll() ([]user.ResponseDTO, error)
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

	var saveDto dto.SaveRequestDTO
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
