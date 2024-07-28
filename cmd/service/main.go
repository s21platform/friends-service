package main

import (
	"github.com/s21platform/friends-service/internal/config"
	db "github.com/s21platform/friends-service/internal/repositore/db"
	"log"
)

func main() {
	//чтение конфига
	cfg := config.MustLoad()
	Data, err := db.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer Data.Connection.Close()

	// добавление grpc сервера

}
