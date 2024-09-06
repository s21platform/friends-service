package consumer_notification_new_user //nolint:revive,stylecheck

import (
	"context"
	"fmt"
	"log"

	"github.com/s21platform/friends-service/internal/repository/db"
	"github.com/s21platform/friends-service/internal/repository/kafka/producer/producer_notification_new_user"

	"github.com/s21platform/friends-service/internal/config"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	Consumer                *kafka.Reader
	NotificationNewPeerProd *producer_notification_new_user.KafkaProducer
	dbRepo                  *db.Repository
}

func New(
	cfg *config.Config,
	prod *producer_notification_new_user.KafkaProducer,
	dbRepo *db.Repository,
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

	return &KafkaConsumer{Consumer: reader,
		NotificationNewPeerProd: prod,
		dbRepo:                  dbRepo}, nil
}

func (kc *KafkaConsumer) Listen() {
	for {
		readMsg := kc.process()
		writeMsg, err := kc.dbRepo.GetUUIDForEmail(readMsg)

		if err != nil {
			fmt.Println("Not work: ", err)
			continue
		}

		err = kc.NotificationNewPeerProd.Process(string(readMsg), writeMsg)

		if err != nil {
			fmt.Println("NewUserProd.process: ", err)
			break
		}

		err = kc.dbRepo.UpdateUserInvite(string(readMsg))

		if err != nil {
			fmt.Println("Not update DB: ", err)
		}
	}
}

func (kc *KafkaConsumer) readMessage() (kafka.Message, error) {
	msg, err := kc.Consumer.ReadMessage(context.Background())

	if err != nil {
		return kafka.Message{}, fmt.Errorf("kc.ReadMessage: %v", err)
	}

	return msg, nil
}

func (kc *KafkaConsumer) process() []byte {
	msg, err := kc.readMessage()

	if err != nil {
		log.Println("Error read message: ", err)
	}

	return msg.Value
}
