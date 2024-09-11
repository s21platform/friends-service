package user_new_peer //nolint:revive,stylecheck

import (
	"context"
	"encoding/json"
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

type FriendRegisterRsvMap struct {
	Email string `json:"email"`
	UUID  string `json:"uuid"`
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
		readMsg, err := kc.readMessage()

		if err != nil {
			fmt.Println("kc.process() ", err)
			continue
		}

		writeMsg, err := kc.dbR.GetUUIDForEmail(readMsg.Email)

		if err != nil {
			fmt.Println("Not work: ", err)
			continue
		}

		err = kc.notificationNewPeerProd.Process(readMsg.Email, writeMsg)

		if err != nil {
			fmt.Println("NewUserProd.process: ", err)
			continue
		}
	}
}

func (kc *KafkaConsumer) readMessage() (FriendRegisterRsvMap, error) {
	var Friend FriendRegisterRsvMap
	msgJson, err := kc.consumer.ReadMessage(context.Background())
	if err != nil {
		return Friend, err
	}

	err = json.Unmarshal(msgJson.Value, &Friend)
	if err != nil {
		return Friend, err
	}

	log.Println("read from topic (email):", Friend.Email)

	return Friend, nil
}
