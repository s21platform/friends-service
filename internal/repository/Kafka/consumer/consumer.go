package consumer

import (
	"context"
	"fmt"

	"github.com/s21platform/friends-service/internal/config"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	Consumer *kafka.Reader
}

func New(cfg *config.Config) (*KafkaConsumer, error) {
	broker := []string{cfg.Kafka.Server}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: broker,
		Topic:   cfg.Kafka.TopicForReading,
	})

	ctx, cansel := context.WithCancel(context.Background())
	defer cansel()

	_, err := reader.ReadMessage(ctx)

	if err != nil {
		return nil, fmt.Errorf("kafka.NewReader: %v", err)
	}

	return &KafkaConsumer{Consumer: reader}, nil
}

func (kc *KafkaConsumer) ReadMessage() (kafka.Message, error) {
	msg, err := kc.Consumer.ReadMessage(context.Background())

	if err != nil {
		return kafka.Message{}, fmt.Errorf("kc.ReadMessage: %v", err)
	}

	return msg, nil
}
