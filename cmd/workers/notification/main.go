package main

import (
	"context"
	"log"

	"github.com/s21platform/friends-service/internal/databus/notification"
	kafkalib "github.com/s21platform/kafka-lib"
	"github.com/s21platform/metrics-lib/pkg"

	"github.com/s21platform/friends-service/internal/config"
	"github.com/s21platform/friends-service/internal/repository/postgres"
)

func main() {
	cfg := config.MustLoad()
	dbRepo, err := postgres.New(cfg)

	if err != nil {
		log.Fatalf("postgres.New: %v", err)
	}

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, "friends", cfg.Platform.Env)
	if err != nil {
		log.Fatalf("faild to connect graphite: %v", err)
	}

	ctx := context.WithValue(context.Background(), config.KeyMetrics, metrics)

	// Consumers
	newFriendConsumer, err := kafkalib.NewConsumer(cfg.Kafka.Server, cfg.Kafka.TopicNewFriend, metrics)
	if err != nil {
		log.Fatalf("failed to create consumer: %v", err)
	}

	// Producers
	notificationNewFriendProducer := kafkalib.NewProducer(cfg.Kafka.Server, cfg.Kafka.NotificationNewFriendTopic)

	// Kafka Handlers
	NewFriendHandler := notification.New(dbRepo, notificationNewFriendProducer)

	// Register Handlers
	newFriendConsumer.RegisterHandler(ctx, NewFriendHandler.Handler)

	<-ctx.Done()
}
