package user_new_peer //nolint:revive,stylecheck

import (
	"context"
	"fmt"
	"github.com/s21platform/friends-service/internal/config"
	"github.com/segmentio/kafka-go"
	"log"
)

type KafkaConsumer struct {
	consumer                *kafka.Reader
	notificationNewPeerProd ProdRepo
	dbR                     DBRepo
}

func New(
	cfg *config.Config,
	prod ProdRepo,
	storage DBRepo,
) (*KafkaConsumer, error) {
	broker := []string{cfg.Kafka.Server}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: broker,
		Topic:   cfg.Kafka.TopicForReading,
		GroupID: "123",
	})

	return &KafkaConsumer{consumer: reader,
		notificationNewPeerProd: prod,
		dbR:                     storage}, nil
}

func (kc *KafkaConsumer) Listen() {
	for {
		readMsg, err := kc.process()

		if err != nil {
			fmt.Println("kc.process() ", err)
			continue
		}

		writeMsg, err := kc.dbR.GetUUIDForEmail(readMsg)

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

	log.Println("read topic: ", msg.Value)

	return msg, nil
}

func (kc *KafkaConsumer) process() ([]byte, error) {
	msg, err := kc.readMessage()

	if err != nil {
		return nil, fmt.Errorf("kc.readMessage: %v", err)
	}

	return msg.Value, nil
}
