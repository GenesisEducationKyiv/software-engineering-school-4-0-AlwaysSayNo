package mail

import (
	"context"
	"encoding/json"
	"fmt"
	myproducer "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/mail/producer"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-AlwaysSayNo/pkg/broker/producer"
	"strconv"
	"time"
)

const (
	CurrencyUpdatedEvent         = "CurrencyUpdated"
	UserSubscribedEvent          = "UserSubscribed"
	UserSubscriptionUpdatedEvent = "UserSubscriptionUpdated"
)

type Event struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type Producer interface {
	Publish(ctx context.Context, body []byte) error
	Close() error
}

type Client struct {
	producer      Producer
	lastCommandID int
}

func NewClient(cnf myproducer.Config) (*Client, error) {
	queueProducer, err := producer.NewProducer(producer.Config(cnf))
	if err != nil {
		return nil, fmt.Errorf("creating consumer: %w", err)
	}

	return &Client{
		producer:      queueProducer,
		lastCommandID: 0,
	}, nil
}

func (c *Client) SendEvent(ctx context.Context, eventType string, data any) error {
	cmd := c.createCommand(eventType, data)

	body, err := c.marshal(cmd)
	if err != nil {
		return fmt.Errorf("marshaling: %w", err)
	}

	if err = c.producer.Publish(ctx, body); err != nil {
		return fmt.Errorf("publishing message: %w", err)
	}

	return nil
}

func (c *Client) createCommand(eventType string, data any) Event {
	c.lastCommandID++
	return Event{
		ID:        strconv.Itoa(c.lastCommandID),
		Type:      eventType,
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      data,
	}
}

func (c *Client) marshal(event Event) ([]byte, error) {
	body, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("marshalling event: %w", err)
	}

	return body, nil
}

func (c *Client) Close() error {
	if err := c.producer.Close(); err != nil {
		return fmt.Errorf("closing producer: %w", err)
	}

	return nil
}
