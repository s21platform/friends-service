package main

import (
	"context"
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
}
