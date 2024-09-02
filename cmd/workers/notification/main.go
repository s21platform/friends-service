package main

import (
	"context"
	"fmt"
	"github.com/s21platform/friends-service/internal/repository/Kafka/consumer"
	"log"
	"time"

	"github.com/s21platform/friends-service/internal/config"
	"github.com/s21platform/friends-service/internal/repository/Kafka/producer"
)

func main() {
	env := config.MustLoad()

	prod, err := producer.New(env)
	if err != nil {
		log.Println("Error create produser: ", err)
	}

	defer prod.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Microsecond)

	defer cancel()

	err = prod.SendMessage(ctx, []byte("Hello, test"))
	if err != nil {
		log.Println("Error sendMessage: ", err)
	}

	reader, err := consumer.New(env)
	defer reader.Close()
	if err != nil {
		log.Println("Error create reader: ", err)
	}

	for {
		msg, err := reader.ReadMessage(10 * time.Microsecond)
		if err != nil {
			log.Fatalf("Error reading message: %v", err)
		}

		fmt.Printf("Message at offset %d from topic %s:\nKey: %s, Value: %s\n",
			msg.Offset, msg.Topic, string(msg.Key), string(msg.Value))
	}
}
