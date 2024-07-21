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
	if err != nil {
		log.Fatal(err)
	}
	defer db.Connection.Close()
}
