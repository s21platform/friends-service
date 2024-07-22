package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/s21platform/friends-service/internal/broker"
	config2 "github.com/s21platform/friends-service/internal/config"
	"time"
)

func main() {
	// Параметры подключения к Kafka
	env := config2.MustLoad()

	cfgMap := &kafka.ConfigMap{
		"bootstrap.servers": env.Kafka.Server,
		"group.id":          env.Kafka.GroupId,
	}

	consumer, err := broker.New(cfgMap)
	if err != nil {
		fmt.Printf("Error creating consumer: %s\n", err)
	}
	defer consumer.Consumer.Close()
	err = consumer.Consumer.SubscribeTopics([]string{env.Kafka.TopicForReading}, nil)
	if err != nil {
		fmt.Printf("Error subscribing: %s\n", err)
	}
	//тут должно вернуться сообщение, которое мы отправим в другой топик
	_, err = consumer.ReadMessage(100 * time.Millisecond)
	if err != nil {
		fmt.Printf("Error reading message: %s\n", err)
	}
}
