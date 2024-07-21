package main

import (
	"context"
	"database/sql"
	friend_proto "github.com/s21platform/friends-proto/friends-proto"
)

type Server struct {
	friend_proto.UnimplementedFriendsServiseServer
	Postgres *sql.DB
}

func (s *Server) SetFriend(ctx context.Context, in *friend_proto.SetFriendsIn) (*friend_proto.SetFriendsOut, error) {
	_, _ = ctx, in

	return &friend_proto.SetFriendsOut{Success: true}, nil
}
