package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func New() (*sql.DB, error) {
	//Connect db
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		"myuser", "mypassword", "mydatabase", "localhost", 5433)

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
