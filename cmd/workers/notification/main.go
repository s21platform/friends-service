package main

import (
	"context"
	"fmt"

	"log"

	"github.com/s21platform/friends-service/internal/config"
	"github.com/s21platform/friends-service/internal/repository/Kafka/consumer"
	"github.com/s21platform/friends-service/internal/repository/Kafka/producer"
)

func main() {
	env := config.MustLoad()
	fmt.Println(env)
	prod, err := producer.New(env)

	if err != nil {
		log.Println("Error create produser: ", err)
	}

	defer prod.Close()

	cons, err := consumer.New(env)

	if err != nil {
		log.Println("Error create consumer: ", err)
	}

	for {
		msg, err := cons.ReadMessage()
		if err != nil {
			log.Println("Error read message: ", err)
		}

		fmt.Println("Hi", msg.Value)

		err = prod.SendMessage(context.Background(), []byte("Hello, test"))
		if err != nil {
			log.Println("Error sendMessage: ", err)
		}
	}
}
