package user

import (
	"fmt"
	"log"
	"time"

	"github.com/s21platform/friends-service/internal/config"
	user_proto "github.com/s21platform/user-proto/user-proto"
	"google.golang.org/grpc"
)

type Handle struct {
	client user_proto.UserServiceClient
}

func connect(cfg *config.Config) (*Handle, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.User.Host, cfg.User.Port))

	if err != nil {
		return nil, fmt.Errorf("grpc.NewClient: %w", err)
	}

	client := user_proto.NewUserServiceClient(conn)

	return &Handle{client: client}, nil
}

func New(cfg *config.Config) (*Handle, error) {
	var err error

	var client *Handle

	for i := 0; i < 5; i++ {
		client, err = connect(cfg)
		if err == nil {
			return client, nil
		}

		log.Println(err)
		time.Sleep(500 * time.Microsecond)
	}

	return nil, err
}
