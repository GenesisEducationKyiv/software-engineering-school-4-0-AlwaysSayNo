package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/broker/consumer"
	"github.com/AlwaysSayNo/genesis-currency-api/email-service/internal/module/broker/consumer/config"
	"log"
	"time"
)

const (
	MailTimeout       = 10 * time.Second
	MailerCommandType = "SendEmails"
)

type Listener func([]byte) error

type Mailer interface {
	SendEmail(ctx context.Context, emails []string, subject, message string) error
}

type Consumer interface {
	Subscribe(listener Listener)
	Listen(stop <-chan struct{})
	Close() error
}

type Command struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	Data      Data   `json:"data"`
}

type Data struct {
	Emails  []string `json:"emails"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

type Client struct {
	stop     chan struct{}
	consumer Consumer
}

func NewClient(cnf config.ConsumerConfig) (*Client, error) {
	emailConsumer, err := consumer.NewConsumer(cnf)
	if err != nil {
		return nil, fmt.Errorf("creating consumer: %w", err)
	}

	stop := make(chan struct{})
	go emailConsumer.Listen(stop)

	return &Client{
		stop:     stop,
		consumer: emailConsumer,
	}, nil
}

func (c *Client) Subscribe(ctx context.Context, mailer Mailer) error {
	c.consumer.Subscribe(func(body []byte) error {
		ctx, cancel := context.WithTimeout(ctx, MailTimeout)
		defer cancel()

		cmd, err := unmarshal(body)
		if err != nil {
			return fmt.Errorf("unmarshaling: %w", err)
		}

		if cmd.Type != MailerCommandType {
			return nil
		}

		log.Printf("Command (id: %d, timestamp: %s)", cmd.ID, cmd.Timestamp)

		return mailer.SendEmail(ctx, cmd.Data.Emails, cmd.Data.Subject, cmd.Data.Body)
	})

	return nil
}

func unmarshal(body []byte) (*Command, error) {
	command := Command{}

	if err := json.Unmarshal(body, &command); err != nil {
		return nil, fmt.Errorf("unmarshalling response body to command: %w", err)
	}

	return &command, nil
}

func (c *Client) Close() error {
	close(c.stop)
	return c.consumer.Close()
}
