package notifier

import (
	"bytes"
	"context"
	"fmt"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/notifier/config"
	sharcurrdto "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/shared/dto/currency"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/shared/dto/user"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/pkg/apperrors"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

type UserGetter interface {
	GetAll() ([]user.ResponseDTO, error)
}

type DatedCurrencyGetter interface {
	GetCachedCurrency(ctx context.Context) (sharcurrdto.CachedCurrency, error)
}

type MailClient interface {
	SendEmail(ctx context.Context, emails []string, subject string, message string) error
	Close() error
}

type CurrencyEmailData struct {
	FromCcy    string
	ToCcy      string
	UpdateDate string
	BuyRate    float64
	SaleRate   float64
}

type EmailNotifier struct {
	userGetter     UserGetter
	currencyGetter DatedCurrencyGetter
	mailClient     MailClient
	cnf            config.EmailServiceConfig
}

func NewEmailNotifier(
	mailClient MailClient,
	datedCurrencyGetter DatedCurrencyGetter,
	userGetter UserGetter,
	cnf config.EmailServiceConfig,
) *EmailNotifier {
	return &EmailNotifier{
		mailClient:     mailClient,
		currencyGetter: datedCurrencyGetter,
		userGetter:     userGetter,
		cnf:            cnf,
	}
}

func (em *EmailNotifier) Notify(ctx context.Context) error {
	log.Println("Start sending emails")
	body, err := em.prepareEmail(ctx)
	if err != nil {
		return err
	}

	if err := em.send(ctx, body); err != nil {
		log.Println("Unsuccessful finish emails sending")
		return fmt.Errorf("sending email: %w", err)
	}

	log.Println("Finish sending emails")
	return nil
}

// prepareEmail is used to prepare an email. Email consists of an email_template and rate information.
// Returns prepared email or error.
func (em *EmailNotifier) prepareEmail(ctx context.Context) (*bytes.Buffer, error) {
	// Get an email_template.
	tmplPath := filepath.Join("pkg", "common", "templates", "email_template.html")

	tmpl, err := os.ReadFile(tmplPath)
	if err != nil {
		return nil, apperrors.NewInvalidStateError("reading the template file", err)
	}

	t, err := template.New("email").Parse(string(tmpl))
	if err != nil {
		return nil, apperrors.NewInvalidStateError("parsing the template file", err)
	}

	currency, err := em.currencyGetter.GetCachedCurrency(ctx)
	if err != nil {
		return nil, err
	}

	// Put currency data into email_template
	var body bytes.Buffer
	err = t.Execute(&body, currency)
	if err != nil {
		return nil, apperrors.NewInvalidStateError("executing template:", err)
	}

	return &body, nil
}

func (em *EmailNotifier) send(ctx context.Context, body *bytes.Buffer) error {
	// Empty users list check.
	users, err := em.userGetter.GetAll()
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

	err = em.mailClient.SendEmail(ctx, emails, em.cnf.EmailSubject, body.String())
	if err != nil {
		return fmt.Errorf("sending emails: %w", err)
	}

	return nil
}
