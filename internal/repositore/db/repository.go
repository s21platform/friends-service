package db

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/s21platform/friends-service/internal/config"
	"log"
	"path/filepath"
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

func (r *Repository) MigrateDB() error {
	driver, err := postgres.WithInstance(r.Connection, &postgres.Config{})
	if err != nil {
		log.Fatal("error getting driver", err)
	}

	migrationsPath, err := filepath.Abs("../../scripts/migrations")
	if err != nil {
		log.Fatal("error getting absolute path", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsPath, "postgres", driver)
	if err != nil {
		log.Fatal("error getting migrate object", err, migrationsPath)
	}

	//Применение миграций
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("error migration process", err)
	}
	return nil
}
