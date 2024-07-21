package kafka

import (
	"database/sql"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/s21platform/friends-service/internal/config"
	"log"
)

var db *sql.DB

func processMessage(cfg *config.Config, email string) {
	// Проверка есть ли uuid в Бд
	rows, err := db.Query("SELECT uuid FROM peer_invite WHERE email = $1", email)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var uuids []string
	for rows.Next() {
		var uuid string
		if err := rows.Scan(&uuid); err != nil {
			log.Fatal(err)
		}
		uuids = append(uuids, uuid)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	// Отправка UUID в другой топик
	for _, uuid := range uuids {
		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &notificationsTopic, Partition: kafka.PartitionAny},
			Value:          []byte(fmt.Sprintf("UUID: %s", uuid)),
		}, nil)
		if err != nil {
			log.Printf("Failed to deliver message: %v", err)
		}
	}
}
