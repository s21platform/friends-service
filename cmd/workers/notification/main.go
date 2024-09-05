package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/s21platform/friends-service/internal/repository/db"

	"github.com/s21platform/friends-service/internal/config"
	"github.com/s21platform/friends-service/internal/repository/Kafka/consumer"
	"github.com/s21platform/friends-service/internal/repository/Kafka/producer"
)

func main() {
	env := config.MustLoad()
	dbRepo, err := db.New(env)

	if err != nil {
		log.Printf("db.New: %v", err)
		os.Exit(1)
	}

	prod, err := producer.New(env)

	if err != nil {
		log.Println("Error create produser: ", err)
	}

	defer prod.Close()

	cons, err := consumer.New(env)

	if err != nil {
		log.Println("Error create consumer: ", err)
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
