package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/mail/producer"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/mail/producer/config"
	"strconv"
	"time"
)

const (
	eventType = "SendEmail"
)

type Command struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	Data      any    `json:"data"`
}

type data struct {
	Emails  []string `json:"emails"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

type Producer interface {
	Publish(ctx context.Context, body []byte) error
	Close() error
}

type Client struct {
	producer      Producer
	lastCommandID int
}

func NewClient(cnf config.ProducerConfig) (*Client, error) {
	emailProducer, err := producer.NewProducer(cnf)
	if err != nil {
		return nil, fmt.Errorf("creating consumer: %w", err)
	}

	return &Client{
		producer:      emailProducer,
		lastCommandID: 0,
	}, nil
}

func (c *Client) SendEmail(
	ctx context.Context,
	emails []string, subject string, message string) error {
	data := data{
		Emails:  emails,
		Subject: subject,
		Body:    message,
	}

	cmd := c.createCommand(data)

	body, err := c.marshal(cmd)
	if err != nil {
		return fmt.Errorf("marshaling: %w", err)
	}

	if err = c.producer.Publish(ctx, body); err != nil {
		return fmt.Errorf("publishing message: %w", err)
	}

	return nil
}

func (c *Client) createCommand(data data) Command {
	c.lastCommandID++
	return Command{
		ID:        strconv.Itoa(c.lastCommandID),
		Type:      eventType,
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      data,
	}
}

func (c *Client) marshal(command Command) ([]byte, error) {
	body, err := json.Marshal(command)
	if err != nil {
		return nil, fmt.Errorf("marshalling command: %w", err)
	}

	return body, nil
}

func (c *Client) Close() error {
	if err := c.producer.Close(); err != nil {
		return fmt.Errorf("closing producer: %w", err)
	}

	return nil
}
