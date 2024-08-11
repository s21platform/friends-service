package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Friends  Service
	Postgres ReadEnvBD
	Kafka    Kafka
}

type Service struct {
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

type Kafka struct {
	TopicForReading string `env:"FRIENDS_SERVICE_NOTIFICATION_KAFKA_TOPIC"`
	TopicForWriting string `env:"NOTIFICATION_SERVICE_FRIENDS_TOPIC"`
	Server          string `env:"KAFKA_SERVER" envDefault:"localhost:9092"`
	GroupId         string `env:"KAFKA_GROUP_ID" envDefault:"test"`
	AutoOffset      string `env:"KAFKA_OFFSET" envDefault:"latest"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("Can not read env variables: %s", err)
	}
	return cfg
}
