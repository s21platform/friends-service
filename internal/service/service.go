package service

import (
	"context"
	friend_proto "github.com/s21platform/friends-proto/friends-proto"
)

type Server struct {
	friend_proto.UnimplementedFriendsServiseServer
	dbR dbRepo
}

func (s *Server) SetFriend(ctx context.Context, in *friend_proto.SetFriendsIn) (*friend_proto.SetFriendsOut, error) {
	res, err := s.dbR.isRowFriendExist(in.Peer_1, in.Peer_2)
	if err != nil || res == false {
		return &friend_proto.SetFriendsOut{Success: false}, nil
	}
	res, err = s.dbR.SetFriend(in.Peer_1, in.Peer_2)
	if err != nil || res == false {
		return &friend_proto.SetFriendsOut{Success: false}, nil
	}
	return &friend_proto.SetFriendsOut{Success: true}, nil
}

func (s *Server) New(repo dbRepo) *Server {
	return &Server{dbR: repo}
}
