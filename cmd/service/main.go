package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/s21platform/friends-service/internal/rpc/user"

	friends "github.com/s21platform/friends-proto/friends-proto"
	"github.com/s21platform/friends-service/internal/config"
	db "github.com/s21platform/friends-service/internal/repository/db"
	"github.com/s21platform/friends-service/internal/service"
	"google.golang.org/grpc"
)

func main() {
	// чтение конфига
	cfg := config.MustLoad()
	dbRepo, err := db.New(cfg)

	if err != nil {
		log.Printf("db.New: %v", err)
		os.Exit(1)
	}
	defer dbRepo.Close()

	// добавление grpc сервера
	thisService := service.New(dbRepo)

	s := grpc.NewServer()
	friends.RegisterFriendsServiceServer(s, thisService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Printf("Cannot listen port: %s; Error: %s", cfg.Service.Port, err)
	}

	_, err = user.New(cfg)
	if err != nil {
		log.Printf("Cannot connect User service: %s", err)
	}

	if err = s.Serve(lis); err != nil {
		log.Printf("Cannot start service: %s; Error: %s", cfg.Service.Port, err)
	}
}
