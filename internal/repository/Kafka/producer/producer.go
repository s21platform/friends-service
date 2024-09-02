package producer

import (
	"context"
	"fmt"

	"github.com/s21platform/friends-service/internal/config"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	Producer *kafka.Writer
	Topic    *config.Kafka
}

func New(cfg *config.Config) (*KafkaProducer, error) {
	if cfg.Kafka.Server == "" {
		return nil, fmt.Errorf("kafka server address is not provided")
	}

	if cfg.Kafka.TopicForWriting == "" {
		return nil, fmt.Errorf("kafka topic is not provided")
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Kafka.Server),
		Topic:        cfg.Kafka.TopicForWriting,
		Balancer:     &kafka.LeastBytes{}, // балансировщик, на данный момент равномерно распределяет сообщения по партициям
		RequiredAcks: kafka.RequireAll,    // подтверждение о том что сообщение доставлено
	}

	return &KafkaProducer{Producer: writer, Topic: &cfg.Kafka}, nil
}

func (kp *KafkaProducer) Close() error {
	return kp.Producer.Close()
}

func (kp *KafkaProducer) sendMessage(ctx context.Context, value []byte) error {
	err := kp.Producer.WriteMessages(ctx, kafka.Message{
		Value: value,
	})

	return err
}

func (kp *KafkaProducer) Process(msg []byte) error {
	err := kp.sendMessage(context.Background(), msg)

	if err != nil {
		return fmt.Errorf("kp.sendMessage: %v", err)
	}

	return nil
}
