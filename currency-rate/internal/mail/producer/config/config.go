package config

import "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/pkg/envs"

type ProducerConfig struct {
	BrokerURI string
	QueueName string
}

func LoadProducerConfig() ProducerConfig {
	return ProducerConfig{
		BrokerURI: envs.Get("CURRENCY_SERVICE_BROKER_URI"),
		QueueName: envs.Get("CURRENCY_SERVICE_QUEUE_NAME"),
	}
}
