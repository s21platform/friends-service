package main

import (
	"context"
	"github.com/s21platform/friends-service/internal/config"
	"github.com/s21platform/friends-service/internal/repository/Kafka/producer"
	"log"
	"time"
)

func main() {
	env := config.MustLoad()

	prod, err := producer.New(env)
	if err != nil {
		log.Println("Error create produser: %s", err)
	}

	defer prod.Close()

	ctx, cansel := context.WithTimeout(context.Background(), 10*time.Microsecond)

	defer cansel()

	err = prod.SendMessage(ctx, []byte("Hello, test"))
	if err != nil {
		log.Println("Error sendMessage: %s", err)
	}
}
