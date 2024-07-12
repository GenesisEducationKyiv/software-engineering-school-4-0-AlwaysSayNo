package config

import "github.com/AlwaysSayNo/genesis-currency-api/email-service/pkg/envs"

type ConsumerConfig struct {
	BrokerURI string
	QueueName string
}

func LoadConsumerConfig() ConsumerConfig {
	return ConsumerConfig{
		BrokerURI: envs.Get("BROKER_URI"),
		QueueName: envs.Get("QUEUE_NAME"),
	}
}
