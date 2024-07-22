package main

import (
	"fmt"
	"github.com/s21platform/friends-service/internal/broker"
	config2 "github.com/s21platform/friends-service/internal/config"
	"log"
)

func main() {
	// Параметры подключения к Kafka
	env := config2.MustLoad()
	fmt.Printf("%+v\n", env)
	//database, err := db.New(env)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//Create consumer
	consumer, err := broker.New(env)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Consumer.Close()

	//Subscribe for topic
	err = consumer.Consumer.SubscribeTopics([]string{env.Kafka.TopicForReading}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %v", err)
	}

	for {
		msg, err := consumer.Consumer.ReadMessage(-1) // -1 означает бесконечное ожидание
		if err == nil {
			fmt.Printf("Recived message: %s\n", msg.Value)
			//ProcessMessage(database, string(msg.Value), producer, env)
		} else {
			log.Printf("Error receiving message: %v\n", err)
		}
	}
}

//func ProcessMessage(r *db.Repository, email string, producer *kafka.Producer, cfg *config2.Config) {
//	// Проверка есть ли uuid в Бд
//	rows, err := r.Connection.Query("SELECT initiator FROM user_invite WHERE invited = $1", email)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer rows.Close()
//
//	var uuids []string
//
//	for rows.Next() {
//		var uuid string
//		if err := rows.Scan(&uuid); err != nil {
//			log.Fatal(err)
//		}
//		uuids = append(uuids, uuid)
//	}
//	if err := rows.Err(); err != nil {
//		log.Fatal(err)
//	}
//	// Отправка UUID в другой топик
//	for _, uuid := range uuids {
//		err = producer.Produce(&kafka.Message{
//			TopicPartition: kafka.TopicPartition{Topic: &cfg.Kafka.TopicForWritting, Partition: kafka.PartitionAny},
//			Value:          []byte(fmt.Sprint(uuid)),
//		}, nil)
//		if err != nil {
//			log.Printf("Failed to deliver message: %v", err)
//		}
//	}
//}
