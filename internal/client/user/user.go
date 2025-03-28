package user

import (
	"context"
	"fmt"
	"log"
	"time"

	logger_lib "github.com/s21platform/logger-lib"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/s21platform/friends-service/internal/config"
	user_proto "github.com/s21platform/user-proto/user-proto"
	"google.golang.org/grpc"
)

type Handle struct {
	client user_proto.UserServiceClient
}

func (h *Handle) IsUserExistByUUID(ctx context.Context, userUUID string) (bool, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("IsUserExistByUUID")
	res, err := h.client.IsUserExistByUUID(ctx, &user_proto.IsUserExistByUUIDIn{Uuid: userUUID})

	if err != nil {
		logger.Error("failed to answer user-service")
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
