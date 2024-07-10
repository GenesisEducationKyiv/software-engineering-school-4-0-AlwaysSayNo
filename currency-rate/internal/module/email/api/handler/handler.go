package handler

import (
	"context"
	"github.com/AlwaysSayNo/genesis-currency-api/common/pkg/apperrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmailSender interface {
	SendEmails(ctx context.Context) error
}

type Handler struct {
	emailSender EmailSender
}

func NewHandler(emailSender EmailSender) *Handler {
	return &Handler{
		emailSender: emailSender,
	}
}

func (h *Handler) SendEmails(ctx *gin.Context) {
	err := h.emailSender.SendEmails(ctx)
	if err != nil {
		apperrors.AttachToCtx(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, "Emails are successfully sent")
}
