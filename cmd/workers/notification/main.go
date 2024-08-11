package main

import (
	"fmt"
	"github.com/s21platform/friends-service/internal/config"
	"github.com/s21platform/friends-service/internal/repository/Kafka/consumer"
	"time"
)

func main() {
	env := config.MustLoad()
	cons, err := consumer.New(env)
	if err != nil {
		fmt.Printf("Error creating consumer: %s\n", err)
	}
	defer cons.Consumer.Close()

	//тут должно вернуться сообщение, которое мы отправим в другой топик
	_, err = cons.ReadMessage(100 * time.Millisecond)
	if err != nil {
		fmt.Printf("Error reading message: %s\n", err)
	}
}
