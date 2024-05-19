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

func NewEmailService(userService *UserService, currencyService *CurrencyService) *EmailService {
	return &EmailService{
		userService,
		currencyService,
	}
}

func (s *EmailService) SendEmails() error {
	tmplPath := filepath.Join("pkg", "common", "templates", "email_template.html")

	tmpl, err := os.ReadFile(tmplPath)
	if err != nil {
		return errors.NewInvalidStateError("Failed to read the file", err)
	}

	t, err := template.New("email").Parse(string(tmpl))
	if err != nil {
		return errors.NewInvalidStateError("Failed to parse the file", err)
	}

	rate := s.currencyService.GetCurrencyInfo()

	var body bytes.Buffer
	err = t.Execute(&body, rate)
	if err != nil {
		return errors.NewInvalidStateError("Failed to execute template:", err)
	}

	return s.send(&body)
}

func (s *EmailService) send(body *bytes.Buffer) error {
	log.Println("Start sending emails")

	smtpHost := viper.Get("SMTP_HOST").(string)
	smtpPort := viper.Get("SMTP_PORT").(string)

	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

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

	log.Println(smtpUser, smtpPassword, smtpHost)
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, to, message)
	log.Println(smtpHost+":"+smtpPort, auth, smtpUser, to)
	if err != nil {
		return errors.NewInvalidStateError("Failed to send email:", err)
	}

	log.Println("Finish sending emails")

	return nil
}
