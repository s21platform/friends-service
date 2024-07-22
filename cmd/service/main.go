package main

import (
	"github.com/s21platform/friends-service/internal/config"
	"github.com/s21platform/friends-service/internal/repository/db"
	"log"
)

func main() {
	//чтение конфига
	cfg := config.MustLoad()
	data, err := db.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Connection.Close()
}
