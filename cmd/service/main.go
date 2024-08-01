package main

import (
	"fmt"
	"github.com/s21platform/friends-service/internal/config"
	db "github.com/s21platform/friends-service/internal/repositore/db"
	"log"
)

func main() {
	//чтение конфига
	cfg := config.MustLoad()
	Data, err := db.New(cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("db.New: %w", err))
	}
	//миграции
	err = Data.MigrateDB()
	if err != nil {
		log.Fatal(fmt.Errorf("Data.MigrateDB: %w", err))
	}
	defer Data.Close()

	// добавление grpc сервера

}
