package consumer

import "github.com/GenesisEducationKyiv/software-engineering-school-4-0-AlwaysSayNo/pkg/envs"

type Config struct {
	BrokerURI string
	QueueName string
}

func LoadConsumerConfig() Config {
	return Config{
		BrokerURI: envs.Get("EMAIL_SERVICE_BROKER_URI"),
		QueueName: envs.Get("EMAIL_SERVICE_QUEUE_NAME"),
	}
}
