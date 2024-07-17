package config

import "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/pkg/envs"

type ProducerConfig struct {
	BrokerURI string
	QueueName string
}

func LoadProducerConfig() ProducerConfig {
	return ProducerConfig{
		BrokerURI: envs.Get("BROKER_URI"),
		QueueName: envs.Get("QUEUE_NAME"),
	}
}
