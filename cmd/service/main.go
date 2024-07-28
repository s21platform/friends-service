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
	log.Printf("postgres log host=%s port=%s user=%s dbname=%s\n", cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Database)
	fmt.Printf("postgres fmt host=%s port=%s user=%s dbname=%s\n", cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Database)
	Data, err := db.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer Data.Connection.Close()

	// добавление grpc сервера

}
