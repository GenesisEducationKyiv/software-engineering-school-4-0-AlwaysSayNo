package decorator

import (
	"context"
	"fmt"
	producerclient "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/mail"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/dto"
	userdto "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/shared/dto/user"
)

type Service interface {
	Save(ctx context.Context, saveRequestDTO dto.SaveRequestDTO) (*userdto.ResponseDTO, error)
	GetAll(ctx context.Context) ([]userdto.ResponseDTO, error)
	ChangeSubscriptionStatus(ctx context.Context, id int, isSubscribed bool) (*userdto.ResponseDTO, error)
}

type ProducerClient interface {
	SendEvent(ctx context.Context, eventType string, data any) error
}

type Decorator struct {
	service        Service
	producerClient ProducerClient
}

func NewServiceDecorator(service Service, producerClient ProducerClient) *Decorator {
	return &Decorator{
		service:        service,
		producerClient: producerClient,
	}
}

type UserSubscribedData struct {
	Email        string `json:"email"`
	IsSubscribed bool   `json:"isSubscribed"`
}

type UserSubscriptionUpdatedData struct {
	Email        string `json:"email"`
	IsSubscribed bool   `json:"isSubscribed"`
}

func (d *Decorator) Save(ctx context.Context, saveRequestDTO dto.SaveRequestDTO) (*userdto.ResponseDTO, error) {
	responseDTO, err := d.service.Save(ctx, saveRequestDTO)
	if err != nil {
		return responseDTO, fmt.Errorf("proxing by decorator: %w", err)
	}

	data := UserSubscribedData{
		Email:        saveRequestDTO.Email,
		IsSubscribed: true,
	}

	//todo SAGA
	if err := d.producerClient.SendEvent(ctx, producerclient.UserSubscribedEvent, data); err != nil {
		return nil, fmt.Errorf("proxing by decorator: %w", err)
	}

	return responseDTO, err
}

func (d *Decorator) GetAll(ctx context.Context) ([]userdto.ResponseDTO, error) {
	return d.service.GetAll(ctx)
}

func (d *Decorator) ChangeSubscriptionStatus(ctx context.Context, id int, isSubscribed bool) (*userdto.ResponseDTO, error) {
	responseDTO, err := d.service.ChangeSubscriptionStatus(ctx, id, isSubscribed)
	if err != nil {
		return responseDTO, fmt.Errorf("proxing by decorator: %w", err)
	}

	data := UserSubscriptionUpdatedData{
		Email:        responseDTO.Email,
		IsSubscribed: true,
	}

	//todo SAGA
	if err := d.producerClient.SendEvent(ctx, producerclient.UserSubscriptionUpdatedEvent, data); err != nil {
		return nil, fmt.Errorf("proxing by decorator: %w", err)
	}

	return responseDTO, err
}
