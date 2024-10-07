package decorator

import (
	"context"
	producerclient "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/mail"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/model"
)

type Repository interface {
	Create(ctx context.Context, user model.User) (*model.User, error)
	GetAll(ctx context.Context) (*[]model.User, error)
	ExistsByEmail(ctx context.Context, email string) bool
}

type ProducerClient interface {
	SendEvent(ctx context.Context, eventType string, data any) error
}

type Decorator struct {
	repository     Repository
	producerClient ProducerClient
}

func NewRepositoryDecorator(repository Repository, producerClient ProducerClient) *Decorator {
	return &Decorator{
		repository:     repository,
		producerClient: producerClient,
	}
}

type UserSubscribedData struct {
	Email        string `json:"email"`
	IsSubscribed bool   `json:"isSubscribed"`
}

func (d *Decorator) Create(ctx context.Context, user model.User) (*model.User, error) {
	data := UserSubscribedData{
		Email:        user.Email,
		IsSubscribed: true,
	}

	//todo SAGA
	if err := d.producerClient.SendEvent(ctx, producerclient.UserSubscribedEvent, data); err != nil {
		return nil, err
	}

	return d.repository.Create(ctx, user)
}

//todo change subscription method

func (d *Decorator) GetAll(ctx context.Context) (*[]model.User, error) {
	return d.repository.GetAll(ctx)
}

func (d *Decorator) ExistsByEmail(ctx context.Context, email string) bool {
	return d.repository.ExistsByEmail(ctx, email)
}
