package handler

import (
	"context"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-AlwaysSayNo/pkg/apperrors"
	"net/http"
	"strconv"

	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/dto"
	userdto "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/shared/dto/user"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	Save(ctx context.Context, user dto.SaveRequestDTO) (*userdto.ResponseDTO, error)
	GetAll(ctx context.Context) ([]userdto.ResponseDTO, error)
	ChangeSubscriptionStatus(ctx context.Context, id int, isSubscribed bool) (*userdto.ResponseDTO, error)
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

	_, err := h.userService.Save(ctx, saveDto)
	if err != nil {
		apperrors.AttachToCtx(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, "E-mail додано")
}

func (h *Handler) FindAll(ctx *gin.Context) {
	result, err := h.userService.GetAll(ctx)
	if err != nil {
		apperrors.AttachToCtx(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, &result)
}

func (h *Handler) ChangeSubscriptionStatus(ctx *gin.Context) {
	isSubscribedStr := ctx.Query("isSubscribed")

	isSubscribed, err := strconv.ParseBool(isSubscribedStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value for parameter isSubscribed"})
		return
	}

	idStr := ctx.Query("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value for parameter id"})
		return
	}

	_, err = h.userService.ChangeSubscriptionStatus(ctx, id, isSubscribed)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	message := "User was successfully subscribed"
	if !isSubscribed {
		message = "User was successfully unsubscribed"
	}

	ctx.JSON(http.StatusOK, gin.H{"message": message})
}
