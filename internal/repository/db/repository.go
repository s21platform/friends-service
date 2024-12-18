package db

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq" // Импортируем данную библиотеку для работы с бд.
	"github.com/s21platform/friends-service/internal/config"
)

type Repository struct {
	connection *sqlx.DB
}

func connect(cfg *config.Config) (*Repository, error) {
	// Connect db
	conStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Database)

	db, err := sqlx.Connect("postgres", conStr)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	return &Repository{db}, nil
}

func (r *Repository) GetWhoFollowsPeer(initiator string) ([]string, error) {
	var result []string
	err := r.connection.Select(&result, "SELECT initiator FROM friends WHERE user_id = $1", initiator)

	if err != nil {
		return nil, fmt.Errorf("r.connection.Select: %v", err)
	}

	return result, nil
}

func (r *Repository) GetPeerFollows(initiator string) ([]string, error) {
	var peers []string
	err := r.connection.Select(&peers, "SELECT user_id FROM friends WHERE initiator = $1", initiator)

	if err != nil {
		return nil, fmt.Errorf("r.connection.Select: %v", err)
	}

	return peers, nil
}

func (r *Repository) IsRowFriendExist(peer1, peer2 string) (bool, error) {
	var res []string
	err := r.connection.Select(&res, "SELECT user_id FROM friends WHERE initiator = $1 AND user_id = $2", peer1, peer2)

	if err != nil {
		return false, fmt.Errorf("r.connection.Select: %v", err)
	} else if len(res) == 0 {
		return false, nil
	}

	return true, nil
}

func (r *Repository) SetFriend(peer1, peer2 string) (bool, error) {
	res, err := r.IsRowFriendExist(peer1, peer2)
	if err != nil || !res {
		return false, fmt.Errorf("r.isRowFriendExist: %v", err)
	}

	if _, err = r.connection.Exec("INSERT INTO friends (initiator, user_id) VALUES ($1, $2)", peer1, peer2); err != nil {
		return false, fmt.Errorf("r.connection.Exec: %v", err)
	}

	return true, nil
}

func (r *Repository) RemoveFriends(peer1, peer2 string) (bool, error) {
	res, err := r.IsRowFriendExist(peer1, peer2)
	if err != nil || res {
		return false, fmt.Errorf("r.isRowFriendExist: %v", err)
	}

	if _, err = r.connection.Exec("DELETE FROM friends WHERE initiator = $1 AND user_id = $2", peer1, peer2); err != nil {
		return false, fmt.Errorf("r.connection.Exec: %v", err)
	}

	return true, nil
}

func (r *Repository) RemoveSubscribe(peer1, peer2 string) error {
	friend, err := r.IsRowFriendExist(peer1, peer2)
	if err != nil {
		return err
	} else if friend {
		return fmt.Errorf("RemoveSubscribe not friend: %s %s", peer1, peer2)
	}

	if _, err = r.connection.Exec("DELETE FROM friends WHERE initiator = $1 AND user_id = $2", peer1, peer2); err != nil {
		return err
	}

	return nil
}

func (r *Repository) isRowInviteExist(row1, row2 string) (bool, error) {
	var res []string
	err := r.connection.Select(&res, "SELECT 1 FROM user_invite WHERE initiator = $1 AND invited = $2", row1, row2)

	if err != nil {
		return false, fmt.Errorf("r.connection.Select: %v", err)
	} else if len(res) == 0 {
		return false, nil
	}

	return true, err
}

func (r *Repository) SetInvitePeer(uuid, email string) error {
	res, err := r.isRowInviteExist(uuid, email)
	if err != nil {
		return err
	} else if !res {
		return fmt.Errorf("invite exist: %s and %s", uuid, email) // можно сделать nil если это не ошибка
	}

	_, err = r.connection.Exec(
		"INSERT INTO user_invite (initiator, invited, is_closed) "+
			"VALUES ($1, $2, false)", uuid, email)

	//добавить USER_INVITE_NOTIFICATION

	return err
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

func (r *Repository) GetUUIDForEmail(email string) ([]string, error) {
	var res []string

	err := r.connection.Select(&res, "SELECT initiator FROM user_invite WHERE invited = $1 AND is_closed != true",
		email)

	if err != nil {
		return nil, fmt.Errorf("r.connection.Select: %v", err)
	}

	return res, nil
}

func (r *Repository) UpdateUserInvite(initiator, invited string) error {
	_, err := r.connection.Exec("UPDATE user_invite SET is_closed=true WHERE initiator=$1 AND invited=$2",
		initiator,
		invited,
	)

	if err != nil {
		return fmt.Errorf("r.connection.Exec UPDATE user_invite: %v", err)
	}

	return nil
}

func (r *Repository) GetCountFriends(uuid string) (int64, int64, error) {
	var subscription, subscribers int64

	err := r.connection.Get(&subscribers, "SELECT count(initiator) FROM friends WHERE user_id = $1", uuid)

	if err != nil {
		return 0, 0, fmt.Errorf("subscription r.connection.Select: %v", err)
	}

	err = r.connection.Get(&subscription, "SELECT count(user_id) FROM friends WHERE initiator = $1", uuid)

	if err != nil {
		return 0, 0, fmt.Errorf("subscribers strconv.ParseInt: %v", err)
	}

	log.Println("subscription, subscribers: ", subscription, subscribers)
	return subscription, subscribers, nil
}
