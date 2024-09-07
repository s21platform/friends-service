package notification_new_user //nolint:revive,stylecheck

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/s21platform/friends-service/internal/config"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	producer *kafka.Writer
	dbR      DBRepo
}

func New(cfg *config.Config, dbR DBRepo) (*KafkaProducer, error) {
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

	return &KafkaProducer{producer: writer, dbR: dbR}, nil
}

func (kp *KafkaProducer) Close() error {
	return kp.producer.Close()
}

func (kp *KafkaProducer) sendMessage(ctx context.Context, email, value string) error {
	jsonMsg := struct {
		Email string `json:"email"`
		UUID  string `json:"uuid"`
	}{
		Email: email,
		UUID:  value,
	}

	msg, err := json.Marshal(jsonMsg)

	if err != nil {
		return fmt.Errorf("json.Marshal: %v", err)
	}

	err = kp.producer.WriteMessages(ctx, kafka.Message{
		Value: msg,
	})
	if err != nil {
		return fmt.Errorf("kp.producer.WriteMessages: %v", err)
	}

	return nil
}

func (kp *KafkaProducer) Process(email string, msgs []string) error {
	for _, val := range msgs {
		err := kp.sendMessage(context.Background(), email, val)

		if err != nil {
			return fmt.Errorf("kp.sendMessage: %v", err)
		}

		err = kp.dbR.UpdateUserInvite(val, email)

		if err != nil {
			return fmt.Errorf("kp.storage.UpdateUserInvite: %v", err)
		}
	}

	return nil
}
