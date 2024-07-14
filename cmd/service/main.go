package main

import (
	"github.com/s21platform/friends-service/internal/config"
	db2 "github.com/s21platform/friends-service/internal/repositore/db"
	"log"
)

func main() {
	//чтение конфига
	cfg := config.MustLoad()
	db, err := db2.New(cfg)
	defer db.Connection.Close()
	if err != nil {
		log.Fatal(err)
	}
}
