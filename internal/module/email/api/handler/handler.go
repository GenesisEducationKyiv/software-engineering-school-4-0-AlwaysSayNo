package handler

import (
	"net/http"

	"genesis-currency-api/pkg/errors"
	"github.com/gin-gonic/gin"
)

type EmailSender interface {
	SendEmails() error
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
	err := h.emailSender.SendEmails()
	if err != nil {
		errors.AttachToCtx(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, "Emails are successfully sent")
}
