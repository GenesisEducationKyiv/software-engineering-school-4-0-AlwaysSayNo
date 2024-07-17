package handler

import (
	"context"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/pkg/apperrors"
	"net/http"

	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/dto"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/shared/dto/user"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	Save(user dto.SaveRequestDTO) (*user.ResponseDTO, error)
	GetAll() ([]user.ResponseDTO, error)
}

type MailNotifier interface {
	Notify(ctx context.Context) error
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
		apperrors.AttachToCtx(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, "E-mail додано")
}

func (h *Handler) FindAll(ctx *gin.Context) {
	result, err := h.userService.GetAll()
	if err != nil {
		apperrors.AttachToCtx(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, &result)
}
