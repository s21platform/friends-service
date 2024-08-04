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

func (s *Server) GetPeerFollow(ctx context.Context, in *friend_proto.GetPeerFollowIn) (*friend_proto.GetPeerFollowOut, error) {
	peersUUID, err := s.dbR.GetPeerFollows(in.Uuid)
	if err != nil {
		return nil, err
	}
	var peers []*friend_proto.Peer
	for _, uuid := range peersUUID {
		peers = append(peers, &friend_proto.Peer{Uuid: uuid})
	}
	return &friend_proto.GetPeerFollowOut{Subscription: peers}, nil
}

func (s *Server) GetWhoFollowPeer(ctx context.Context, in *friend_proto.GetWhoFollowPeerIn) (*friend_proto.GetWhoFollowPeerOut, error) {
	peerUuid, err := s.dbR.GetWhoFollowsPeer(in.Uuid)
	if err != nil {
		return nil, err
	}
	var peers []*friend_proto.Peer
	for _, uuid := range peerUuid {
		peers = append(peers, &friend_proto.Peer{Uuid: uuid})
	}
	return &friend_proto.GetWhoFollowPeerOut{Subscribers: peers}, nil
}
