package service

import (
	"context"

	friend_proto "github.com/s21platform/friends-proto/friends-proto"
)

type Server struct {
	friend_proto.UnimplementedFriendsServiceServer
	dbR DBRepo
}

func (s *Server) SetFriends(
	ctx context.Context, in *friend_proto.SetFriendsIn,
) (*friend_proto.SetFriendsOut, error) {
	_ = ctx
	res, err := s.dbR.SetFriend(in.Peer_1, in.Peer_2)

	if err != nil || !res {
		return nil, err
	}

	return &friend_proto.SetFriendsOut{Success: true}, nil
}

func (s *Server) RemoveSubscribe(
	ctx context.Context, in *friend_proto.RemoveSubscribeIn,
) (*friend_proto.RemoveSubscribeOut, error) {
	_ = ctx
	err := s.dbR.RemoveSubscribe(in.Peer_1, in.Peer_2)

	return &friend_proto.RemoveSubscribeOut{}, err
}

func New(repo DBRepo) *Server {
	return &Server{dbR: repo}
}

func (s *Server) GetPeerFollow(
	ctx context.Context, in *friend_proto.GetPeerFollowIn,
) (*friend_proto.GetPeerFollowOut, error) {
	_ = ctx
	peersUUID, err := s.dbR.GetPeerFollows(in.Uuid)

	if err != nil {
		return nil, err
	}

	peers := make([]*friend_proto.Peer, 0)

	for _, uuid := range peersUUID {
		peers = append(peers, &friend_proto.Peer{Uuid: uuid})
	}

	return &friend_proto.GetPeerFollowOut{Subscription: peers}, nil
}

func (s *Server) GetWhoFollowPeer(
	ctx context.Context, in *friend_proto.GetWhoFollowPeerIn,
) (*friend_proto.GetWhoFollowPeerOut, error) {
	_ = ctx
	peerUUID, err := s.dbR.GetWhoFollowsPeer(in.Uuid)

	if err != nil {
		return nil, err
	}

	peers := make([]*friend_proto.Peer, 0)

	for _, uuid := range peerUUID {
		peers = append(peers, &friend_proto.Peer{Uuid: uuid})
	}

	return &friend_proto.GetWhoFollowPeerOut{Subscribers: peers}, nil
}

func (s *Server) GetInvitePeer(
	ctx context.Context, in *friend_proto.GetInvitePeerIn,
) (*friend_proto.GetInvitePeerOut, error) {
	_ = ctx
	err := s.dbR.GetInvitePeer(in.Uuid, in.Email)

	return &friend_proto.GetInvitePeerOut{}, err
}
