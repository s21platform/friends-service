package broker

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/s21platform/friends-service/internal/config"
)

type KafkaConsumer struct {
	Config   *kafka.ConfigMap
	Consumer *kafka.Consumer
}

func New(cfg *config.Config) (*KafkaConsumer, error) {
	kafkaCfg := &kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.Server,
		"group.id":          cfg.Kafka.GroupId,
		"auto.offset.reset": cfg.Kafka.AutoOffset,
	}
	consumer, err := kafka.NewConsumer(kafkaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	return &KafkaConsumer{
		Config:   kafkaCfg,
		Consumer: consumer,
	}, nil
}
