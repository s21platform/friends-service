package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // Импортируем данную библиотеку для работы с бд.
	"github.com/s21platform/friends-service/internal/config"
)

type Repository struct {
	connection *sql.DB
}

func connect(cfg *config.Config) (*Repository, error) {
	// Connect db
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Host, cfg.Postgres.Port)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// Сhecking connection db
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	return &Repository{db}, nil
}

func (r *Repository) GetWhoFollowsPeer(initiator string) ([]string, error) {
	row, err := r.connection.Query("SELECT initiator FROM friends WHERE user_id = $1", initiator)
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

	if err = row.Err(); err != nil {
		log.Println("rows.Err: ", err)
		return nil, err
	}

	return result, nil
}

func (r *Repository) GetPeerFollows(initiator string) ([]string, error) {
	row, err := r.connection.Query("SELECT user_id FROM friends WHERE initiator = $1", initiator)
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

	if err = row.Err(); err != nil {
		log.Println("rows.Err: ", err)
		return nil, err
	}

	return peers, nil
}

func (r *Repository) SetFriend(peer1, peer2 string) (bool, error) {
	res, err := r.isRowFriendExist(peer1, peer2)
	if err != nil || res {
		return false, err
	}

	if _, err = r.connection.Exec("INSERT INTO friends (initiator, user_id) VALUES ($1, $2)", peer1, peer2); err != nil {
		return false, err
	}

	return true, nil
}

func (r *Repository) RemoveSubscribe(peer1, peer2 string) error {
	friend, err := r.isRowFriendExist(peer1, peer2)
	if err != nil {
		return err
	} else if !friend {
		return fmt.Errorf("RemoveSubscribe not friend: %s %s", peer1, peer2)
	}

	if _, err = r.connection.Exec("DELETE FROM friends WHERE initiator = $1 AND user_id = $2", peer1, peer2); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetInvitePeer(uuid, email string) error {
	res, err := r.isRowInviteExist(uuid, email)
	if err != nil {
		return err
	} else if res {
		return fmt.Errorf("invite exist: %s and %s", uuid, email) // можно сделать nil если это не ошибка
	}

	_, err = r.connection.Exec(
		"INSERT INTO user_invite (initiator, invited, is_closed) "+
			"VALUES ($1, $2, false)", uuid, email)

	return err
}

func (r *Repository) isRowInviteExist(row1, row2 string) (bool, error) {
	row, err := r.connection.Query("SELECT * FROM user_invite WHERE initiator = $1 AND invited = $2", row1, row2)
	if err == nil {
		if !row.Next() {
			return false, err
		}
	}
	defer row.Close()

	if err = row.Err(); err != nil {
		log.Println("rows.Err: ", err)
		return false, err
	}

	return true, nil
}

func (r *Repository) isRowFriendExist(peer1, peer2 string) (bool, error) {
	row, err := r.connection.Query("SELECT user_id FROM friends WHERE initiator = $1 AND user_id = $2", peer1, peer2)
	if err == nil {
		if !row.Next() {
			return false, nil
		}
	}
	defer row.Close()

	if err = row.Err(); err != nil {
		log.Println("rows.Err: ", err)
		return false, err
	}

	return true, err
}

func (r *Repository) Close() {
	r.connection.Close()
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
