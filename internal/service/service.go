package service

import (
	"context"
	"fmt"

	"github.com/s21platform/friends-service/internal/config"

	friend_proto "github.com/s21platform/friends-proto/friends-proto"
	_ "google.golang.org/grpc"
)

type Server struct {
	friend_proto.UnimplementedFriendsServiceServer
	dbR DBRepo
}

func (s *Server) SetFriends(
	ctx context.Context, in *friend_proto.SetFriendsIn,
) (*friend_proto.SetFriendsOut, error) {
	userID := ctx.Value(config.KeyUUID).(string)
	if userID == "" {
		return nil, fmt.Errorf("uuid not found in context")
	}
	res, err := s.dbR.SetFriend(userID, in.Peer)

	if err != nil || !res {
		return nil, err
	}

	return &friend_proto.SetFriendsOut{Success: true}, nil
}

func (s *Server) RemoveFriends(
	ctx context.Context, in *friend_proto.RemoveFriendsIn,
) (*friend_proto.RemoveFriendsOut, error) {
	userID := ctx.Value(config.KeyUUID)
	if userID == nil {
		return nil, fmt.Errorf("uuid not found in context")
	}
	res, err := s.dbR.RemoveFriends(userID.(string), in.Peer)

	if err != nil || !res {
		return nil, err
	}

	return &friend_proto.RemoveFriendsOut{Success: true}, err
}

func (s *Server) RemoveSubscribe(
	ctx context.Context, in *friend_proto.RemoveSubscribeIn,
) (*friend_proto.RemoveSubscribeOut, error) {
	userID := ctx.Value(config.KeyUUID)
	if userID == nil {
		return nil, fmt.Errorf("uuid not found in context")
	}
	err := s.dbR.RemoveSubscribe(userID.(string), in.Peer)

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

func (s *Server) SetInvitePeer(
	ctx context.Context, in *friend_proto.SetInvitePeerIn,
) (*friend_proto.SetInvitePeerOut, error) {
	userID := ctx.Value(config.KeyUUID)
	if userID == nil {
		return nil, fmt.Errorf("uuid not found in context")
	}
	err := s.dbR.SetInvitePeer(userID.(string), in.Email)

	// или тут добавить USER_INVITE_NOTIFICATION

	return &friend_proto.SetInvitePeerOut{}, err
}

func (s *Server) GetCountFriends(ctx context.Context, in *friend_proto.EmptyFriends) (*friend_proto.GetCountFriendsOut, error) {
	userID := ctx.Value(config.KeyUUID)
	if userID == nil {
		return nil, fmt.Errorf("uuid not found in context")
	}
	subscription, subscribers, err := s.dbR.GetCountFriends(userID.(string))
	if err != nil {
		return nil, err
	}

	return &friend_proto.GetCountFriendsOut{
		Subscription: subscription,
		Subscribers:  subscribers,
	}, nil
}
