package email

import (
	"bytes"
	"context"
	"fmt"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/dto"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/service"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/service/email/config"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

type UserGetter interface {
	GetAllSubscribed(ctx context.Context) ([]dto.UserResponseDTO, error)
}

type CurrencyGetter interface {
	FindLatest(ctx context.Context) (*dto.CurrencyDTO, error)
}

type Service struct {
	userGetter     UserGetter
	currencyGetter CurrencyGetter
	cnf            config.EmailServiceConfig
}

func NewEmailService(userGetter UserGetter,
	currencyGetter CurrencyGetter,
	cnf config.EmailServiceConfig) *Service {
	return &Service{
		userGetter:     userGetter,
		currencyGetter: currencyGetter,
		cnf:            cnf,
	}
}

func (es *Service) SendCurrencyPriceEmails(ctx context.Context, mailTransport service.Mailer) error {
	log.Println("Start sending currency price emails")
	body, err := es.prepareCurrencyEmail(ctx)
	if err != nil {
		return err
	}

	if err := es.sendCurrencyEmail(ctx, body, mailTransport); err != nil {
		log.Println("Unsuccessful finish currency price emails sending")
		return fmt.Errorf("sending email: %w", err)
	}

	log.Println("Finish sending currency price emails")
	return nil
}

func (es *Service) prepareCurrencyEmail(ctx context.Context) (*bytes.Buffer, error) {
	// Get an email_template.
	tmplPath := filepath.Join("pkg", "templates", "currency_email_template.html")

	tmpl, err := os.ReadFile(tmplPath)
	if err != nil {
		return nil, fmt.Errorf("reading the template file: %w", err)
	}

	t, err := template.New("currency_email").Parse(string(tmpl))
	if err != nil {
		return nil, fmt.Errorf("parsing the template file: %w", err)
	}

	currency, err := es.currencyGetter.FindLatest(ctx)
	if err != nil {
		return nil, err
	}

	// Put currency data into email_template
	var body bytes.Buffer
	err = t.Execute(&body, currency)
	if err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}

	return &body, nil
}

func (es *Service) sendCurrencyEmail(ctx context.Context, body *bytes.Buffer, mailTransport service.Mailer) error {
	// Empty users list check.
	users, err := es.userGetter.GetAllSubscribed(ctx)
	if err != nil {
		return fmt.Errorf("fetching users: %w", err)
	}
	if len(users) == 0 {
		log.Println("Emails list is empty")
		return nil
	}

	emails := make([]string, 0, len(users))
	for _, u := range users {
		emails = append(emails, u.Email)
	}

	err = mailTransport.SendEmail(ctx, emails, es.cnf.EmailSubject, body.String())
	if err != nil {
		return fmt.Errorf("sending emails: %w", err)
	}

	return nil
}
