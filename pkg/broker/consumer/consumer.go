package consumer

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
)

var listenerMutex = sync.Mutex{}

type Listener func([]byte) error

type Consumer struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	messages  <-chan amqp.Delivery
	listeners []Listener //todo how to do this properly
	cnf       Config
}

func NewConsumer(cnf Config) (*Consumer, error) {
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

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, fmt.Errorf("setting QoS: %w", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, fmt.Errorf("delivery creating: %w", err)
	}

	return &Consumer{
		conn:      conn,
		channel:   ch,
		messages:  msgs,
		listeners: make([]Listener, 0),
		cnf:       cnf,
	}, nil
}

func (c *Consumer) Subscribe(listener Listener) {
	log.Printf("Start subscribing new listener")

	listenerMutex.Lock()
	defer listenerMutex.Unlock()

	c.listeners = append(c.listeners, listener)

	log.Printf("Finish subscribing new listener. Length: %d", len(c.listeners))
}

func (c *Consumer) Listen(stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			return
		case msg, ok := <-c.messages:
			if !ok {
				return
			}

			log.Println("Message received")
			c.handleMessage(msg)
		}
	}
}

func (c *Consumer) handleMessage(msg amqp.Delivery) {
	listenerMutex.Lock()
	defer listenerMutex.Unlock()

	for _, listener := range c.listeners {
		err := listener(msg.Body)
		if err != nil {
			log.Printf("Cannot send message: %v", err)
		}
	}
}

func (c *Consumer) Close() error {
	if err := c.channel.Close(); err != nil {
		return fmt.Errorf("closing channel: %w", err)
	}
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("closing connection: %w", err)
	}

	log.Println("Stop listening")
	return nil
}
