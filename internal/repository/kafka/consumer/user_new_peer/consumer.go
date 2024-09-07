package user_new_peer //nolint:revive,stylecheck

import (
	"context"
	"fmt"

	"github.com/s21platform/friends-service/internal/config"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	consumer                *kafka.Reader
	notificationNewPeerProd ProdRepo
	storage                 Storage
}

func New(
	cfg *config.Config,
	prod ProdRepo,
	storage Storage,
) (*KafkaConsumer, error) {
	broker := []string{cfg.Kafka.Server}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: broker,
		Topic:   cfg.Kafka.TopicForReading,
		GroupID: "123",
	})

	ctx, cansel := context.WithCancel(context.Background())
	defer cansel()

	_, err := reader.ReadMessage(ctx)

	if err != nil {
		return nil, fmt.Errorf("kafka.NewReader: %v", err)
	}

	return &KafkaConsumer{consumer: reader,
		notificationNewPeerProd: prod,
		storage:                 storage}, nil
}

func (kc *KafkaConsumer) Listen() {
	for {
		readMsg, err := kc.process()

		if err != nil {
			fmt.Println("kc.process() ", err)
			continue
		}

		writeMsg, err := kc.storage.GetUUIDForEmail(readMsg)

		if err != nil {
			fmt.Println("Not work: ", err)
			continue
		}

		err = kc.notificationNewPeerProd.Process(string(readMsg), writeMsg)

		if err != nil {
			fmt.Println("NewUserProd.process: ", err)
			continue
		}
	}
}

func (kc *KafkaConsumer) readMessage() (kafka.Message, error) {
	msg, err := kc.consumer.ReadMessage(context.Background())

	if err != nil {
		return kafka.Message{}, fmt.Errorf("kc.ReadMessage: %v", err)
	}

	return msg, nil
}

func (kc *KafkaConsumer) process() ([]byte, error) {
	msg, err := kc.readMessage()

	if err != nil {
		return nil, fmt.Errorf("kc.readMessage: %v", err)
	}

	return msg.Value, nil
}
