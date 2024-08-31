package consumer

import (
	"context"
	"github.com/s21platform/friends-service/internal/config"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	Reader *kafka.Reader
}

func New(cfg *config.Config) (*KafkaConsumer, error) {
	brokerList := []string{cfg.Kafka.Server}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokerList,
		Topic:    cfg.Kafka.TopicForReading,
		GroupID:  cfg.Kafka.GroupID,
		MinBytes: 10e+6, // 10MB
		MaxBytes: 10e+7, // 20MB
		Offset:   0,
		MaxWait:  1 * time.Second,
	})

	return &KafkaConsumer{
		Reader: reader,
	}, nil
}

func (kc *KafkaConsumer) ReadMessage(timeout time.Duration) (*kafka.Message, error) {
	startTime := time.Now()
	for {
		msg, err := kc.Reader.FetchMessage(context.Background())
		if err != nil {
			return nil, err
		}

		if startTime.Add(timeout).Before(time.Now()) {
			return &msg, nil
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func (kc *KafkaConsumer) Close() error {
	return kc.Reader.Close()
}

func (kc *KafkaConsumer) CommitMessages(msg *kafka.Message) error {
}
