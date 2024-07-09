package handler

import (
	"context"
	"net/http"

	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/pkg/errors"
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
		errors.AttachToCtx(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, "Emails are successfully sent")
}
