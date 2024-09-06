package main

import (
	"fmt"
	"github.com/s21platform/friends-service/internal/repository/kafka/consumer/consumer_notification_new_user"
	"github.com/s21platform/friends-service/internal/repository/kafka/produser/producer_notification_new_user"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/s21platform/friends-service/internal/repository/db"

	"github.com/s21platform/friends-service/internal/config"
)

func main() {
	env := config.MustLoad()
	dbRepo, err := db.New(env)

	if err != nil {
		log.Printf("db.New: %v", err)
		os.Exit(1)
	}

	prod, err := producer_notification_new_user.New(env)

	if err != nil {
		log.Println("Error create produser: ", err)
	}

	defer prod.Close()

	cons, err := consumer_notification_new_user.New(env)

	if err != nil {
		log.Println("Error create consumer_notification_new_user: ", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			readMsg := cons.Process()
			writeMsg, err := dbRepo.GetUUIDForEmail(readMsg)

			if err != nil {
				fmt.Println("Not work: ", err)
				continue
			}

			err = prod.Process(string(readMsg), writeMsg)

			if err != nil {
				fmt.Println("prod.Process: ", err)
				break
			}

			err = dbRepo.UpdateUserInvite(string(readMsg))

			if err != nil {
				fmt.Println("Not update DB: ", err)
			}
		}
	}()

	<-done
	log.Println("Shutting down gracefully...")
}
