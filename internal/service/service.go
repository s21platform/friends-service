package service

import (
	"context"
	"fmt"

	"github.com/s21platform/friends-service/internal/config"

	friend_proto "github.com/s21platform/friends-proto/friends-proto"
	logger_lib "github.com/s21platform/logger-lib"

	_ "google.golang.org/grpc"
)

type Server struct {
	friend_proto.UnimplementedFriendsServiceServer
	dbR DBRepo
}

func (s *Server) SetFriends(
	ctx context.Context, in *friend_proto.SetFriendsIn,
) (*friend_proto.SetFriendsOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("SetFriends")
	userID := ctx.Value(config.KeyUUID).(string)
	if userID == "" {
		logger.Error("failed to not found UUID in context")
		return nil, fmt.Errorf("uuid not found in context")
	}
	res, err := s.dbR.SetFriend(userID, in.Peer)

	if err != nil {
		logger.Error("failed to SetFriend from BD")
		return nil, err
	}

	if !res {
		logger.Info("user already in friends")
		return &friend_proto.SetFriendsOut{Success: false}, nil
	}

	return &friend_proto.SetFriendsOut{Success: res}, nil
}

func (s *Server) RemoveFriends(
	ctx context.Context, in *friend_proto.RemoveFriendsIn,
) (*friend_proto.RemoveFriendsOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("RemoveFriends")
	userID := ctx.Value(config.KeyUUID).(string)
	if userID == "" {
		logger.Error("failed to not found UUID in context")
		return nil, fmt.Errorf("uuid not found in context")
	}
	res, err := s.dbR.RemoveFriends(userID, in.Peer)

	if err != nil {
		logger.Error("failed to RemoveFriend from BD")
		return nil, err
	}

	return &friend_proto.RemoveFriendsOut{Success: res}, err
}

func (s *Server) RemoveSubscribe(
	ctx context.Context, in *friend_proto.RemoveSubscribeIn,
) (*friend_proto.RemoveSubscribeOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("RemoveSubscribe")
	userID := ctx.Value(config.KeyUUID).(string)
	if userID == "" {
		logger.Error("failed to not found UUID in context")
		return nil, fmt.Errorf("uuid not found in context")
	}
	err := s.dbR.RemoveSubscribe(userID, in.Peer)

	if err != nil {
		logger.Error("failed to RemoveSubscribe from BD")
	}

	return &friend_proto.RemoveSubscribeOut{}, err
}

func New(repo DBRepo) *Server {
	return &Server{dbR: repo}
}

func (s *Server) GetPeerFollow(
	ctx context.Context, in *friend_proto.GetPeerFollowIn,
) (*friend_proto.GetPeerFollowOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetPeerFollow")
	peersUUID, err := s.dbR.GetPeerFollows(in.Uuid)

	if err != nil {
		logger.Error("failed to GetPeerFollow from BD")
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
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetWhoFollowPeer")
	peerUUID, err := s.dbR.GetWhoFollowsPeer(in.Uuid)

	if err != nil {
		logger.Error("failed to GetWhoFollowPeer from BD")
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
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("SetInvitePeer")
	userID := ctx.Value(config.KeyUUID).(string)
	if userID == "" {
		logger.Error("failed to not found UUID in context")
		return nil, fmt.Errorf("uuid not found in context")
	}
	err := s.dbR.SetInvitePeer(userID, in.Email)

	if err != nil {
		logger.Error("failed to SetInvitePeer from BD")
	}

	// или тут добавить USER_INVITE_NOTIFICATION

	return &friend_proto.SetInvitePeerOut{}, err
}

func (s *Server) GetCountFriends(ctx context.Context, in *friend_proto.EmptyFriends) (*friend_proto.GetCountFriendsOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetCountFriends")
	userIDValue := ctx.Value(config.KeyUUID)
	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		logger.Error("failed to not found UUID in context")
		return nil, fmt.Errorf("uuid not found in context")
	}
	subscription, subscribers, err := s.dbR.GetCountFriends(userID)
	if err != nil {
		logger.Error("failed to GetCountFriends from BD")
		return nil, err
	}

	return &friend_proto.GetCountFriendsOut{
		Subscription: subscription,
		Subscribers:  subscribers,
	}, nil
}

func (s *Server) IsFriendExist(ctx context.Context, in *friend_proto.IsFriendExistIn) (*friend_proto.IsFriendExistOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("IsFriendExist")
	userIDValue := ctx.Value(config.KeyUUID)
	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		logger.Error("failed to not found UUID in context")
		return nil, fmt.Errorf("uuid not found in context")
	}
	res, err := s.dbR.IsRowFriendExist(userID, in.Peer)
	if err != nil {
		logger.Error("failed to IsFriendExist from BD")
		return &friend_proto.IsFriendExistOut{Success: false}, err
	}
	return &friend_proto.IsFriendExistOut{Success: res}, nil
}
