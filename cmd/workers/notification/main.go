package main

import (
	"fmt"
	"github.com/s21platform/friends-service/internal/broker"
	"github.com/s21platform/friends-service/internal/config"
	"time"
)

func main() {
	env := config.MustLoad()
	consumer, err := broker.New(env)
	if err != nil {
		fmt.Printf("Error creating consumer: %s\n", err)
	}
	defer consumer.Consumer.Close()

	//тут должно вернуться сообщение, которое мы отправим в другой топик
	_, err = consumer.ReadMessage(100 * time.Millisecond)
	if err != nil {
		fmt.Printf("Error reading message: %s\n", err)
	}
}
