package main

import (
	"log"
	"os"

	"github.com/s21platform/friends-service/internal/repository/kafka/producer/notification_new_user"

	"github.com/s21platform/friends-service/internal/config"
	"github.com/s21platform/friends-service/internal/repository/db"
	"github.com/s21platform/friends-service/internal/repository/kafka/consumer/user_new_peer"
)

func main() {
	cfg := config.MustLoad()
	dbRepo, err := db.New(cfg)

	if err != nil {
		log.Printf("db.New: %v", err)
		os.Exit(1)
	}

	NewUserProd, err := notification_new_user.New(cfg)

	if err != nil {
		log.Println("Error create producer: ", err)
	}

	defer NewUserProd.Close()

	NewUserCons, err := user_new_peer.New(cfg, NewUserProd, dbRepo)

	if err != nil {
		log.Println("Error create user_new_peer: ", err)
	}

	NewUserCons.Listen()
}
