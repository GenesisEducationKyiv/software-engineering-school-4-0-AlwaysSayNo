package service

import (
	"genesis-currency-api/pkg/request"
)

type EmailService struct {
	emails []request.Email
}

func New() *EmailService {
	return &EmailService{}
}

func (service *EmailService) Save(email request.Email) request.Email {
	service.emails = append(service.emails, email)
	return email
}

func (service *EmailService) FindAll() []request.Email {
	return service.emails
}
