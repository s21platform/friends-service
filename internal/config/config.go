package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Friends  Friends
	Postgres ReadEnvBD
}

type Friends struct {
	Port string `env:"FRIENDS_SERVICE_PORT"`
	Host string `env:"FRIENDS_SERVICE_HOST"`
}

type ReadEnvBD struct {
	User     string `env:"FRIENDS_SERVICE_POSTGRES_USER"`
	Password string `env:"FRIENDS_SERVICE_POSTGRES_PASSWORD"`
	Database string `env:"FRIENDS_SERVICE_POSTGRES_DB"`
	Host     string `env:"FRIENDS_SERVICE_POSTGRES_HOST"`
	Port     string `env:"FRIENDS_SERVICE_POSTGRES_PORT"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("Can not read env variables: %s", err)
	}
	return cfg
}
