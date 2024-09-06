package main

import (
	"log"
	"os"

	"github.com/s21platform/friends-service/internal/repository/db"
	"github.com/s21platform/friends-service/internal/repository/kafka/consumer/consumer_notification_new_user"
	"github.com/s21platform/friends-service/internal/repository/kafka/producer/producer_notification_new_user"

	"github.com/s21platform/friends-service/internal/config"
)

func main() {
	cfg := config.MustLoad()
	dbRepo, err := db.New(cfg)

	if err != nil {
		log.Printf("db.New: %v", err)
		os.Exit(1)
	}

	NewUserProd, err := producer_notification_new_user.New(cfg)

	if err != nil {
		log.Println("Error create producer: ", err)
	}

	defer NewUserProd.Close()

	NewUserCons, err := consumer_notification_new_user.New(cfg, NewUserProd, dbRepo)

	if err != nil {
		log.Println("Error create consumer_notification_new_user: ", err)
	}

	NewUserCons.Listen()
}
