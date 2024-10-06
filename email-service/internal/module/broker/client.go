package broker

import (
	"context"
	"encoding/json"
	"fmt"
	myconsumer "github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/broker/consumer"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/mail/dto"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-AlwaysSayNo/pkg/broker/consumer"
	"log"
	"time"
)

const (
	Timeout                      = 10 * time.Second
	CurrencyUpdatedEvent         = "CurrencyUpdated"
	SubscribeUserSubscribedEvent = "UserSubscribed"
	UserSubscriptionUpdatedEvent = "UserSubscriptionUpdated"
)

type Mailer interface {
	SendEmail(ctx context.Context, emails []string, subject, message string) error
}

type CurrencyService interface {
	Save(ctx context.Context, currencyAddDTO dto.CurrencyAddDTO) error
}

type UserService interface {
	Save(ctx context.Context, userSaveDTO dto.UserSaveDTO) error
	ChangeUserSubscriptionStatus(ctx context.Context, email string, isSubscribed bool) error
}

type Consumer interface {
	Subscribe(listener consumer.Listener)
	Listen(stop <-chan struct{})
	Close() error
}

type Event struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Timestamp string          `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

type CurrencyUpdateData struct {
	Number float64 `json:"number"`
	Date   string  `json:"data"`
}

type UserSubscribedData struct {
	Email        string `json:"email"`
	IsSubscribed bool   `json:"isSubscribed"`
}

type UserSubscriptionUpdatedData struct {
	Email        string `json:"email"`
	IsSubscribed bool   `json:"isSubscribed"`
}

type Client struct {
	stop          chan struct{}
	queueConsumer Consumer
}

func NewClient(cnf myconsumer.Config) (*Client, error) {
	queueConsumer, err := consumer.NewConsumer(consumer.Config(cnf))
	if err != nil {
		return nil, fmt.Errorf("creating consumer: %w", err)
	}

	stop := make(chan struct{})
	go queueConsumer.Listen(stop)

	return &Client{
		stop:          stop,
		queueConsumer: queueConsumer,
	}, nil
}

func (c *Client) SubscribeCurrencyUpdateEvent(ctx context.Context, service CurrencyService) error {
	c.queueConsumer.Subscribe(func(body []byte) error {
		ctx, cancel := context.WithTimeout(ctx, Timeout)
		defer cancel()

		event, err := unmarshalEvent(body)
		if err != nil {
			return fmt.Errorf("unmarshalling CurrencyUpdatedEvent: %w", err)
		}

		if event.Type != CurrencyUpdatedEvent {
			return nil
		}

		log.Printf("CurrencyUpdatedEvent (id: %s, timestamp: %s)", event.ID, event.Timestamp)

		var data CurrencyUpdateData
		if err := json.Unmarshal(event.Data, &data); err != nil {
			return fmt.Errorf("unmarshalling CurrencyUpdateData: %w", err)
		}

		if err := service.Save(ctx, dto.CurrencyUpdateDataToCurrencyAddDTO(&data)); err != nil {
			return fmt.Errorf("saving CurrencyUpdateData: %w", err)
		}

		return nil
	})

	return nil
}

func (c *Client) SubscribeUserSubscribedEvent(ctx context.Context, service UserService) error {
	c.queueConsumer.Subscribe(func(body []byte) error {
		ctx, cancel := context.WithTimeout(ctx, Timeout)
		defer cancel()

		event, err := unmarshalEvent(body)
		if err != nil {
			return fmt.Errorf("unmarshalling SubscribeUserSubscribedEvent: %w", err)
		}

		if event.Type != SubscribeUserSubscribedEvent {
			return nil
		}

		log.Printf("SubscribeUserSubscribedEvent (id: %s, timestamp: %s)", event.ID, event.Timestamp)

		var data UserSubscribedData
		if err := json.Unmarshal(event.Data, &data); err != nil {
			return fmt.Errorf("unmarshalling UserSubscribedData: %w", err)
		}

		if err := service.Save(ctx, dto.UserSubscribedDataToUserSaveDTO(&data)); err != nil {
			return fmt.Errorf("saving UserSubscribedData: %w", err)
		}

		return nil
	})

	return nil
}

func (c *Client) SubscribeUserSubscriptionUpdatedEvent(ctx context.Context, service UserService) error {
	c.queueConsumer.Subscribe(func(body []byte) error {
		ctx, cancel := context.WithTimeout(ctx, Timeout)
		defer cancel()

		event, err := unmarshalEvent(body)
		if err != nil {
			return fmt.Errorf("unmarshalling UserSubscriptionUpdatedEvent: %w", err)
		}

		if event.Type != UserSubscriptionUpdatedEvent {
			return nil
		}

		log.Printf("UserSubscriptionUpdatedEvent (id: %s, timestamp: %s)", event.ID, event.Timestamp)

		var data UserSubscriptionUpdatedData
		if err := json.Unmarshal(event.Data, &data); err != nil {
			return fmt.Errorf("unmarshalling UserSubscriptionUpdatedData: %w", err)
		}

		if err := service.ChangeUserSubscriptionStatus(ctx, data.Email, data.IsSubscribed); err != nil {
			return fmt.Errorf("changing user's subscription: %w", err)
		}

		return nil
	})

	return nil
}

func unmarshalEvent(body []byte) (*Event, error) {
	event := Event{}

	if err := json.Unmarshal(body, &event); err != nil {
		return nil, fmt.Errorf("unmarshalling response body to event: %w", err)
	}

	return &event, nil
}

func (c *Client) Close() error {
	close(c.stop)
	return c.queueConsumer.Close()
}
