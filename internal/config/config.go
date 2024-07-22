package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Friends  Friends
	Kafka    Kafka
	Postgres ReadEnvBD
}

type ReadEnvBD struct {
	User     string `env:"POSTGRES_USER" envDefault:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" envDefault:"postgres"`
	Database string `env:"POSTGRES_DB" envDefault:"postgres"`
	Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
}

type Friends struct {
	Port string `env:"FRIENDS_SERVICE_PORT" envDefault:"8080"`
	Host string `env:"FRIENDS_SERVICE_HOST" envDefault:"localhost"`
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
