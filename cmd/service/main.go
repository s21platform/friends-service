package main

import (
	"fmt"
	friend_proto "github.com/s21platform/friends-proto/friends-proto"
	"github.com/s21platform/friends-service/internal/config"
	db "github.com/s21platform/friends-service/internal/repositore/db"
	"github.com/s21platform/friends-service/internal/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	//чтение конфига
	cfg := config.MustLoad()
	dbRepo, err := db.New(cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("db.New: %w", err))
	}
	defer dbRepo.Close()

	//добавление grpc сервера
	thisService := service.New(dbRepo)

	s := grpc.NewServer()
	friend_proto.RegisterFriendsServiseServer(s, thisService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Friends.Port))
	if err != nil {
		log.Fatalf("Cannnot listen port: %s; Error: %s", cfg.Friends.Port, err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Cannnot start service: %s; Error: %s", cfg.Friends.Port, err)
	}
}
