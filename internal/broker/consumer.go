package broker

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"time"
)

type KafkaConsumer struct {
	Consumer *kafka.Consumer
}

func New(cfg *kafka.ConfigMap) (*KafkaConsumer, error) {
	consumer, err := kafka.NewConsumer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
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
