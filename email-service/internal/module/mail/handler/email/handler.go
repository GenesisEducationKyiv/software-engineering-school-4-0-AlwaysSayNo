package email

import (
	"context"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/service"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-AlwaysSayNo/pkg/apperrors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Sender interface {
	SendCurrencyPriceEmails(ctx context.Context, mailTransport service.Mailer) error
}

type Handler struct {
	sender        Sender
	mailTransport service.Mailer
}

func NewHandler(sender Sender, mailTransport *service.Mailer) *Handler {
	return &Handler{
		sender:        sender,
		mailTransport: *mailTransport,
	}
}

func (h *Handler) SendCurrencyPriceEmails(ctx *gin.Context) {
	err := h.sender.SendCurrencyPriceEmails(ctx, h.mailTransport)
	if err != nil {
		apperrors.AttachToCtx(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, "Emails are successfully sent")
}
