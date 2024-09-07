package main

import (
	"log"

	notificationnewuser "github.com/s21platform/friends-service/internal/repository/kafka/producer/notification_new_user"

	"github.com/s21platform/friends-service/internal/config"
	"github.com/s21platform/friends-service/internal/repository/db"
	usernewpeer "github.com/s21platform/friends-service/internal/repository/kafka/consumer/user_new_peer"
)

func main() {
	cfg := config.MustLoad()
	dbRepo, err := db.New(cfg)

	if err != nil {
		log.Fatalf("db.New: %v", err)
	}

	NewUserProd, err := notificationnewuser.New(cfg, dbRepo)

	if err != nil {
		_ = NewUserProd.Close()

		log.Fatalf("Error create producer: %v", err)
	}

	defer NewUserProd.Close()

	NewUserCons, err := usernewpeer.New(cfg, NewUserProd, dbRepo)

	if err != nil {
		_ = NewUserProd.Close()

		log.Fatalf("Error create user_new_peer: %v", err) //nolint:gocritic
	}

	NewUserCons.Listen()
}
