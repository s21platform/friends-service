package new_friend

import (
	"context"
	"encoding/json"
	"log"

	"github.com/s21platform/friends-proto/friends-proto/new_friend_register"
	"github.com/s21platform/friends-service/internal/config"
	"github.com/s21platform/metrics-lib/pkg"
)

// FIXME Перенести определение этой структуры в notification-service

type Message struct {
	Email string `json:"email"`
	UUID  string `json:"uuid"`
}

type Handler struct {
	dbR DBRepo
	nnf NotificationNewFriend
}

func New(dbR DBRepo, nnf NotificationNewFriend) *Handler {
	return &Handler{dbR: dbR, nnf: nnf}
}

func convertMessage(bMessage []byte, target interface{}) error {
	err := json.Unmarshal(bMessage, target)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) Handler(ctx context.Context, in []byte) {
	m := pkg.FromContext(ctx, config.KeyMetrics)

	var msg new_friend_register.NewFriendRegister
	err := convertMessage(in, &msg)
	if err != nil {
		m.Increment("new_friend.error")
		log.Printf("failed to convert message: %v", err)
		return
	}

	uuids, err := h.dbR.GetUUIDForEmail(msg.Email)
	if err != nil {
		m.Increment("new_friend.error")
		log.Printf("failed to get uuid for email: %v", err)
		return
	}

	for _, uuid := range uuids {
		msg := Message{
			Email: msg.Email,
			UUID:  uuid,
		}
		err := h.nnf.ProduceMessage(msg)
		if err != nil {
			m.Increment("new_friend.error")
			log.Printf("failed to produce message for %s: %v", uuid, err)
		}
	}
}
