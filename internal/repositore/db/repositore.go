package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/s21platform/friends-service/internal/config"
)

func New(cfg *config.Config) (*sql.DB, error) {
	//Connect db
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.BD.User, cfg.BD.Password, cfg.BD.Database, cfg.BD.Host, cfg.BD.Port)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		fmt.Println("error connect: ", err)
		return nil, err
	}

	//Сhecking connection db
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
