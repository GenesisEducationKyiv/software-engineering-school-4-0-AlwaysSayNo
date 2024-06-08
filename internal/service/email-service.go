package service

import (
	"bytes"
	"fmt"
	"genesis-currency-api/pkg/errors"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
)

type CurrencyEmailData struct {
	FromCcy    string
	ToCcy      string
	UpdateDate string
	BuyRate    float64
	SaleRate   float64
}

type EmailService struct {
	userService     *UserService
	currencyService *CurrencyService
}

// NewEmailService is a factory function for EmailService
func NewEmailService(userService *UserService, currencyService *CurrencyService) *EmailService {
	return &EmailService{
		userService,
		currencyService,
	}
}

// SendEmails is used to send a currency update email to all subscribed users.
// It sends full CurrencyInfoDto information (buy, sale rates) compared to /api/rate ena-point.
// Returns error in case of occurrence.
func (s *EmailService) SendEmails() error {
	log.Println("Start sending emails")
	body, err := s.prepareEmail()
	if err != nil {
		return err
	}

	return s.send(body)
}

// prepareEmail is used to prepare an email. Email consists of an email_template and rate information.
// Returns prepared email or error.
func (s *EmailService) prepareEmail() (*bytes.Buffer, error) {
	// Get an email_template.
	tmplPath := filepath.Join("pkg", "common", "templates", "email_template.html")

	tmpl, err := os.ReadFile(tmplPath)
	if err != nil {
		return nil, errors.NewInvalidStateError("Failed to read the file", err)
	}

	t, err := template.New("email").Parse(string(tmpl))
	if err != nil {
		return nil, errors.NewInvalidStateError("Failed to parse the file", err)
	}

	rate := s.currencyService.GetCurrencyInfo()

	// Put rate data into email_template
	var body bytes.Buffer
	err = t.Execute(&body, rate)
	if err != nil {
		return nil, errors.NewInvalidStateError("Failed to execute template:", err)
	}

	return &body, nil
}

// send sends emails to users using the standard library.
// If the list of users is empty, it will return an error.
// Returns error in case of occurrence.
func (s *EmailService) send(body *bytes.Buffer) error {

	smtpHost := viper.Get("SMTP_HOST").(string)
	smtpPort := viper.Get("SMTP_PORT").(string)

	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	// Empty users list check.
	users, err := s.userService.GetAll()
	if len(users) == 0 {
		return errors.NewInvalidStateError("Emails list is empty", err)
	}

	var to []string
	for _, u := range users {
		to = append(to, u.Email)
	}

	subject := viper.Get("EMAIL_SUBJECT").(string)
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n"

	message := []byte(fmt.Sprintf("%s\r\n%s\r\n%s", subject, mime, body.String()))
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, to, message)
	if err != nil {
		return errors.NewInvalidStateError("Failed to send email:", err)
	}

	log.Println("Finish sending emails")

	return nil
}
