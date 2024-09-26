package config

import "github.com/GenesisEducationKyiv/software-engineering-school-4-0-AlwaysSayNo/pkg/envs"

type ConsumerConfig struct {
	BrokerURI string
	QueueName string
}

func LoadConsumerConfig() ConsumerConfig {
	return ConsumerConfig{
		BrokerURI: envs.Get("EMAIL_SERVICE_BROKER_URI"),
		QueueName: envs.Get("EMAIL_SERVICE_QUEUE_NAME"),
	}
}
