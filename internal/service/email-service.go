package service

import (
	"bytes"
	"fmt"
	"genesis-currency-api/pkg/dto"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"path/filepath"

	"genesis-currency-api/pkg/config"

	"genesis-currency-api/pkg/errors"
)

type EmailServiceInterface interface {
	SendEmails() error
}

type UserGetter interface {
	GetAll() ([]dto.UserResponseDTO, error)
}

type DatedCurrencyGetter interface {
	GetCachedCurrency() (dto.CachedCurrency, error)
}

type CurrencyEmailData struct {
	FromCcy    string
	ToCcy      string
	UpdateDate string
	BuyRate    float64
	SaleRate   float64
}

type EmailService struct {
	userGetter     UserGetter
	currencyGetter DatedCurrencyGetter
	cnf            config.EmailServiceConfig
}

// NewEmailService is a factory function for EmailService
func NewEmailService(userGetter UserGetter,
	currencyGetter DatedCurrencyGetter,
	cnf config.EmailServiceConfig,
) *EmailService {
	return &EmailService{
		userGetter:     userGetter,
		currencyGetter: currencyGetter,
		cnf:            cnf,
	}
}

// SendEmails is used to send a currency update email to all subscribed users.
// It sends information (rate, update date).
// Returns error in case of occurrence.
func (s *EmailService) SendEmails() error {
	log.Println("Start sending emails")
	body, err := s.prepareEmail()
	if err != nil {
		return err
	}

	if err := s.send(body); err != nil {
		log.Println("Unsuccessful finish emails sending")
		return fmt.Errorf("sending email body: %w", err)
	}

	log.Println("Finish sending emails")
	return nil
}

// prepareEmail is used to prepare an email. Email consists of an email_template and rate information.
// Returns prepared email or error.
func (s *EmailService) prepareEmail() (*bytes.Buffer, error) {
	// Get an email_template.
	tmplPath := filepath.Join("pkg", "common", "templates", "email_template.html")

	tmpl, err := os.ReadFile(tmplPath)
	if err != nil {
		return nil, errors.NewInvalidStateError("reading the template file", err)
	}

	t, err := template.New("email").Parse(string(tmpl))
	if err != nil {
		return nil, errors.NewInvalidStateError("parsing the template file", err)
	}

	currency, err := s.currencyGetter.GetCachedCurrency()
	if err != nil {
		return nil, err
	}

	// Put currency data into email_template
	var body bytes.Buffer
	err = t.Execute(&body, currency)
	if err != nil {
		return nil, errors.NewInvalidStateError("executing template:", err)
	}

	return &body, nil
}

// send sends emails to users using the standard library.
// If the list of users is empty, it will return an error.
// Returns error in case of occurrence.
func (s *EmailService) send(body *bytes.Buffer) error {
	// Empty users list check.
	users, err := s.userGetter.GetAll()
	if len(users) == 0 {
		return errors.NewInvalidStateError("emails list is empty", err)
	}

	to := make([]string, 0, len(users))
	for _, u := range users {
		to = append(to, u.Email)
	}

	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n"

	message := []byte(fmt.Sprintf("%s\r\n%s\r\n%s", s.cnf.EmailSubject, mime, body.String()))
	auth := smtp.PlainAuth("", s.cnf.SMTPUser, s.cnf.SMTPPassword, s.cnf.SMTPHost)

	err = smtp.SendMail(s.cnf.SMTPHost+":"+s.cnf.SMTPPort, auth, s.cnf.SMTPUser, to, message)
	if err != nil {
		return errors.NewInvalidStateError("while sending email:", err)
	}

	log.Println("Finish sending emails")

	return nil
}
