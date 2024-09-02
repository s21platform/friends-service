package consumer

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/s21platform/friends-service/internal/config"
)

type KafkaConsumer struct {
	Consumer *kafka.Consumer
}

func New(cfg *config.Config) (*KafkaConsumer, error) {
	cfgMap := &kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.Server,
		"group.id":          cfg.Kafka.GroupID,
	}

	consumer, err := kafka.NewConsumer(cfgMap)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	err = consumer.SubscribeTopics([]string{cfg.Kafka.TopicForReading}, nil)
	if err != nil {
		fmt.Printf("Error subscribing: %s\n", err)
	}

	return &KafkaConsumer{
		Consumer: consumer,
	}, nil
}

func (kc *KafkaConsumer) ReadMessage(timeout time.Duration) (**kafka.Message, error) {
	// Получаем текущее время
	startTime := time.Now()

	// Цикл ожидания сообщения
	for {
		// Пытаемся получить сообщение с таймаутом
		msg, err := kc.Consumer.ReadMessage(-1)
		if err == nil && msg.Timestamp.Before(startTime.Add(timeout)) {
			// Сообщение получено в пределах таймаута
			return &msg, nil
		} else if err != nil {
			// Произошла ошибка при чтении сообщения
			return nil, err
		}
		// Если сообщение не получено в пределах таймаута, ждем некоторое время перед повторным запросом
		time.Sleep(100 * time.Millisecond)
	}
}
