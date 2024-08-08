package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/s21platform/friends-service/internal/config"
	"log"
	"time"
)

type Repository struct {
	сonnection *sql.DB
}

func connect(cfg *config.Config) (*Repository, error) {
	//Connect db
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Host, cfg.Postgres.Port)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	//Сhecking connection db
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}
	return &Repository{db}, nil
}

func (r *Repository) GetWhoFollowsPeer(initiator string) ([]string, error) {
	row, err := r.сonnection.Query("SELECT initiator FROM friends WHERE user_id = $1", initiator)
	if err != nil {
		log.Println("connection err: ", err)
		return nil, err
	}
	defer row.Close()
	var result []string
	for row.Next() {
		var resStr string
		if err := row.Scan(&resStr); err != nil {
			log.Println(err)
			return nil, err
		}
		result = append(result, resStr)
	}
	return result, nil
}

func (r *Repository) GetPeerFollows(initiator string) ([]string, error) {
	row, err := r.сonnection.Query("SELECT user_id FROM friends WHERE initiator = $1", initiator)
	if err != nil {
		log.Println("connection err: ", err)
		return nil, err
	}
	defer row.Close()
	var peers []string
	for row.Next() {
		var peer string
		if err := row.Scan(&peer); err != nil {
			log.Println("read peer err: ", err)
		}
		peers = append(peers, peer)
	}
	return peers, nil
}

func (r *Repository) SetFriend(peer_1, peer_2 string) (bool, error) {
	res, err := r.isRowFriendExist(peer_1, peer_2)
	if err != nil || res == true {
		return false, err
	}
	_, err = r.сonnection.Exec("INSERT INTO friends (initiator, user_id) VALUES ($1, $2)", peer_1, peer_2)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Repository) isRowFriendExist(peer_1, peer_2 string) (bool, error) {
	row, err := r.сonnection.Query("SELECT user_id FROM friends WHERE initiator = $1 AND user_id = $2", peer_1, peer_2)
	if err == nil {
		if !row.Next() {
			return false, nil
		}
	}
	defer row.Close()
	return true, err
}

func (r *Repository) Close() {
	r.сonnection.Close()
}

func New(cfg *config.Config) (*Repository, error) {
	var err error
	var repo *Repository
	for i := 0; i < 5; i++ {
		repo, err = connect(cfg)
		if err == nil {
			return repo, nil
		}
		log.Println(err)
		time.Sleep(500 * time.Millisecond)
	}
	return nil, err
}
