package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/s21platform/friends-service/internal/infra"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"

	friends "github.com/s21platform/friends-proto/friends-proto"
	"github.com/s21platform/friends-service/internal/client/user"
	"github.com/s21platform/friends-service/internal/config"
	db "github.com/s21platform/friends-service/internal/repository/postgres"
	"github.com/s21platform/friends-service/internal/service"
	"google.golang.org/grpc"
)

func main() {
	// чтение конфига
	cfg := config.MustLoad()
	dbRepo, err := db.New(cfg)

	if err != nil {
		log.Fatalf("failed to postgres.New: %v", err)
		os.Exit(1)
	}
	defer dbRepo.Close()

	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, cfg.Service.Name, cfg.Platform.Env)

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, "friends", cfg.Platform.Env)
	if err != nil {
		log.Fatalln("failed to create metrics:", err)
	}

	// добавление grpc сервера
	thisService := service.New(dbRepo)

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.UnaryInterceptor,
			infra.MetricsInterceptor(metrics),
		),
		grpc.ChainUnaryInterceptor(infra.Logger(logger)),
	)
	friends.RegisterFriendsServiceServer(s, thisService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Printf("failed to cannot listen port: %s; Error: %s", cfg.Service.Port, err)
	}

	_, err = user.New(cfg)
	if err != nil {
		log.Printf("failed to cannot connect User service: %s", err)
	}

	if err = s.Serve(lis); err != nil {
		log.Printf("failed to cannot start service: %s; Error: %s", cfg.Service.Port, err)
	}
}
