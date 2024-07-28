package db

import (
	"database/sql"
	"fmt"
	"log"

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
		log.Println("error connect: ", err)
		return nil, err
	}

	//Сhecking connection db
	if err := db.Ping(); err != nil {
		log.Println("error ping: ", err)
		return nil, err
	}
	return &Repository{db}, nil
}

func (r *Repository) SetFriend(peer_1, peer_2 string) (bool, error) {
	_, err := r.Connection.Exec("INSERT INTO friends (peer_1, peer_2) VALUES ($1, $2)", peer_1, peer_2)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Repository) isRowFriendExist(peer_1, peer_2 string) (bool, error) {
	row, err := r.Connection.Query("SELECT peer_2 FROM friends WHERE $1 AND $2", peer_1, peer_2)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
	}
	defer row.Close()
	return false, err
}
