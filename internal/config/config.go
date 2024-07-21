package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Friends Friends
	Kafka   Kafka
}

type Friends struct {
	Port string `env:"FRIENDS_SERVICE_PORT"`
	Host string `env:"FRIENDS_SERVICE_HOST"`
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
