package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type key string

const KeyMetrics = key("metrics")

type Config struct {
	Service  Service
	Postgres ReadEnvBD
	Kafka    Kafka
	User     User
	Metrics  Metrics
	Platform Platform
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
	TopicNewFriend             string `env:"USER_FRIENDS_REGISTER"`
	NotificationNewFriendTopic string `env:"FRIENDS_EMAIL_INVITE"`
	Server                     string `env:"KAFKA_SERVER"`
	GroupID                    string `env:"KAFKA_GROUP_ID" envDefault:"test"`
	AutoOffset                 string `env:"KAFKA_OFFSET" envDefault:"latest"`
}

type User struct {
	Host string `env:"USER_SERVICE_HOST"`
	Port string `env:"USER_SERVICE_PORT"`
}

type Metrics struct {
	Host string `env:"GRAFANA_HOST"`
	Port int    `env:"GRAFANA_PORT"`
}

type Platform struct {
	Env string `env:"ENV"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)

	if err != nil {
		log.Fatalf("Can not read env variables: %s", err)
	}

	return cfg
}
