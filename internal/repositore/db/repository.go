package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/s21platform/friends-service/internal/config"
)

type Repository struct {
	Connection *sql.DB
}

func New(cfg *config.Config) (*Repository, error) {
	//Connect db
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Host, cfg.Postgres.Port)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		fmt.Println("error connect: ", err)
		return nil, err
	}

	//Сhecking connection db
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Repository{db}, nil
}
