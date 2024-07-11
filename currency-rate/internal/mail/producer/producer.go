package producer

import (
	"context"
	"fmt"
	"github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/mail/producer/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	cnf     config.ProducerConfig
	queue   amqp.Queue
}

func NewProducer(cnf config.ProducerConfig) (*Producer, error) {
	conn, err := amqp.Dial(cnf.BrokerURI)
	if err != nil {
		return nil, fmt.Errorf("dialing amqp: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("getting channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		cnf.QueueName, // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("declaring queue: %w", err)
	}

	return &Producer{
		conn:    conn,
		channel: ch,
		cnf:     cnf,
		queue:   q,
	}, nil
}

func (p *Producer) Publish(ctx context.Context, body []byte) error {
	err := p.channel.PublishWithContext(ctx,
		"",           // exchange
		p.queue.Name, // routing key
		false,        // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})

	if err != nil {
		return fmt.Errorf("publishing with context: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	if err := p.channel.Close(); err != nil {
		return fmt.Errorf("closing channel: %w", err)
	}
	if err := p.conn.Close(); err != nil {
		return fmt.Errorf("closing connection: %w", err)
	}

	log.Println("Stop publishing")
	return nil
}
