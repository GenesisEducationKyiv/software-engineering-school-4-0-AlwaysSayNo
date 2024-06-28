package util

import (
	"net/http"

	"genesis-currency-api/pkg/dto"

	"genesis-currency-api/pkg/errors"
	"github.com/gin-gonic/gin"
)

type UserGetter interface {
	GetAll() ([]dto.UserResponseDTO, error)
}

type EmailSender interface {
	SendEmails() error
}

type Handler struct {
	userGetter  UserGetter
	emailSender EmailSender
}

func NewHandler(userGetter UserGetter,
	emailSender EmailSender) *Handler {
	return &Handler{
		userGetter:  userGetter,
		emailSender: emailSender,
	}
}

func (h *Handler) FindAll(ctx *gin.Context) {
	result, err := h.userGetter.GetAll()
	if err != nil {
		errors.AttachToCtx(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, &result)
}

func (h *Handler) SendEmails(ctx *gin.Context) {
	err := h.emailSender.SendEmails()
	if err != nil {
		errors.AttachToCtx(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, "Emails are successfully sent")
}

// RegisterRoutes registers routes for passed Handler.
func RegisterRoutes(r *gin.Engine, handler Handler) {
	routes := r.Group("/api/v1/util")
	routes.GET("/users", handler.FindAll)
	routes.POST("/emails/send", handler.SendEmails)
}
