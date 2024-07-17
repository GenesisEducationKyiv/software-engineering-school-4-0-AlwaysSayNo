package handler

import (
	"context"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/pkg/apperrors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MailNotifier interface {
	Notify(ctx context.Context) error
}

type Handler struct {
	mailNotifier MailNotifier
}

func NewHandler(emailNotifier MailNotifier) *Handler {
	return &Handler{
		mailNotifier: emailNotifier,
	}
}

func (h *Handler) SendEmails(ctx *gin.Context) {
	err := h.mailNotifier.Notify(ctx)
	if err != nil {
		apperrors.AttachToCtx(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, "Emails are successfully sent")
}
