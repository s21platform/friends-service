package user

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/s21platform/friends-service/internal/config"
	user_proto "github.com/s21platform/user-proto/user-proto"
	"google.golang.org/grpc"
)

type Handle struct {
	client user_proto.UserServiceClient
}

func (h *Handle) GetUserByLogin(ctx context.Context, email string) (string, bool, error) {
	userUUID, err := h.client.GetUserByLogin(ctx, &user_proto.GetUserByLoginIn{Login: email})

	if err != nil {
		return "", false, fmt.Errorf("error client.GetUserByLogin: %w", err)
	}

	return userUUID.Uuid, userUUID.IsNewUser, nil
}

func (h *Handle) IsUserExistByUUID(ctx context.Context, userUUID string) (bool, error) {
	res, err := h.client.IsUserExistByUUID(ctx, &user_proto.IsUserExistByUUIDIn{Uuid: userUUID})

	if err != nil {
		return false, fmt.Errorf("error client.IsUserExistByUUID: %w", err)
	}

	return res.IsExist, nil
}

func connect(cfg *config.Config) (*Handle, error) {
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.User.Host, cfg.User.Port), opts)

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
