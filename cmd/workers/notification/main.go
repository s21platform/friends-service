package main

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/s21platform/friends-service/internal/config"
	"log"
)

func main() {
	//читаем конфиг
	env := config.MustLoad()
	//подключаем бд

	cfg := &kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
	}
	consumer, err := kafka.NewConsumer(cfg)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()
	producer, err := kafka.NewProducer(cfg)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()
	//for {
	//
	//}
}

//var db *sql.DB
