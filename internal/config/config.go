package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Friends  Friends
	Postgres ReadEnvBD
	Kafka    Kafka
}

type Friends struct {
	Port string `env:"FRIENDS_SERVICE_PORT"`
	Host string `env:"FRIENDS_SERVICE_HOST"`
}

type ReadEnvBD struct {
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Database string `env:"POSTGRES_DB"`
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
}

type Kafka struct {
	KafkaTopic        string `env:"FRIENDS_SERVICE_NOTIFICATION_KAFKA_TOPIC"`
	NotificationKafka string `env:"NOTIFICATION_SERVICE_FRIENDS_TOPIC"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("Can not read env variables: %s", err)
	}
	return cfg
}
