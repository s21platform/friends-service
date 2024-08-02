package service

import (
	"context"
	friend_proto "github.com/s21platform/friends-proto/friends-proto"
)

type Server struct {
	friend_proto.UnimplementedFriendsServiseServer
	dbR DbRepo
}

func (s *Server) SetFriends(ctx context.Context, in *friend_proto.SetFriendsIn) (*friend_proto.SetFriendsOut, error) {
	res, err := s.dbR.SetFriend(in.Peer_1, in.Peer_2)
	if err != nil || res == false {
		return &friend_proto.SetFriendsOut{Success: false}, nil
	}
	return &friend_proto.SetFriendsOut{Success: true}, nil
}

func New(repo DbRepo) *Server {
	return &Server{dbR: repo}
}
