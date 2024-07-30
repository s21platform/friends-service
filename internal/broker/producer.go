package broker

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/s21platform/friends-service/internal/config"
)

type KafkaProducer struct {
	Producer *kafka.Producer
	Topic    *config.Kafka
}

func New(cfg *config.Config) (*KafkaProducer, error) {
	cfgMap := &kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.Server,
	}

	producer, err := kafka.NewProducer(cfgMap)
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{
		Producer: producer,
	}, nil
}

func (kp *KafkaProducer) SendMassage(msg string) error {
	sendChan := make(chan kafka.Event)
	err := kp.Producer.Produce(&kafka.Message{TopicPartition: kafka.TopicPartition{
		Topic: &kp.Topic.TopicForWriting, Partition: kafka.PartitionAny,
	}, Value: []byte(msg)}, sendChan)
	if err != nil {
		return err
	}

	for e := range kp.Producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				return fmt.Errorf("%v", ev.TopicPartition.Error)
			} else {
				return nil
			}
		}
	}
	return nil
}
