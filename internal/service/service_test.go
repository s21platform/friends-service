package service

import (
	"context"
	"errors"
	"testing"

	logger_lib "github.com/s21platform/logger-lib"

	"github.com/s21platform/friends-service/internal/config"

	"github.com/docker/distribution/uuid"
	"github.com/golang/mock/gomock"
	friends_proto "github.com/s21platform/friends-proto/friends-proto"
	"github.com/stretchr/testify/assert"
)

func TestServer_GetPeerFollow(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDBRepo := NewMockDBRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

	t.Run("should_ok_with_UUID", func(t *testing.T) {
		userUUID := uuid.Generate()
		followersUUID := []string{
			uuid.Generate().String(),
			uuid.Generate().String(),
		}
		mockDBRepo.EXPECT().GetPeerFollows(userUUID.String()).Return(followersUUID, nil)
		mockLogger.EXPECT().AddFuncName("GetPeerFollow")
		ctx := context.WithValue(ctx, config.KeyLogger, mockLogger)

		s := New(mockDBRepo)
		res, err := s.GetPeerFollow(ctx, &friends_proto.GetPeerFollowIn{Uuid: userUUID.String()})
		assert.NoError(t, err)
		assert.Equal(t, res, &friends_proto.GetPeerFollowOut{Subscription: []*friends_proto.Peer{
			{Uuid: followersUUID[0]},
			{Uuid: followersUUID[1]},
		}})
	})

	t.Run("should_repo_err", func(t *testing.T) {
		userUUID := uuid.Generate()
		repoErr := errors.New("test")

		mockDBRepo.EXPECT().GetPeerFollows(userUUID.String()).Return(nil, repoErr)
		mockLogger.EXPECT().AddFuncName("GetPeerFollow")
		mockLogger.EXPECT().Error("failed to GetPeerFollow from BD")
		ctx := context.WithValue(ctx, config.KeyLogger, mockLogger)

		s := New(mockDBRepo)
		_, err := s.GetPeerFollow(ctx, &friends_proto.GetPeerFollowIn{Uuid: userUUID.String()})
		assert.Error(t, err, repoErr)
	})
}

func TestServer_GetWhoFollowPeer(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDBRepo := NewMockDBRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

	t.Run("should_ok_with_UUID", func(t *testing.T) {
		userUUID := uuid.Generate()
		followersUUID := []string{
			uuid.Generate().String(),
			uuid.Generate().String(),
		}
		mockDBRepo.EXPECT().GetWhoFollowsPeer(userUUID.String()).Return(followersUUID, nil)
		mockLogger.EXPECT().AddFuncName("GetWhoFollowPeer")
		ctx := context.WithValue(ctx, config.KeyLogger, mockLogger)

		s := New(mockDBRepo)
		res, err := s.GetWhoFollowPeer(ctx, &friends_proto.GetWhoFollowPeerIn{Uuid: userUUID.String()})
		assert.NoError(t, err)
		assert.Equal(t, res, &friends_proto.GetWhoFollowPeerOut{Subscribers: []*friends_proto.Peer{
			{Uuid: followersUUID[0]},
			{Uuid: followersUUID[1]},
		}})
	})

	t.Run("should_repo_err", func(t *testing.T) {
		userUUID := uuid.Generate()
		repoErr := errors.New("test")

		mockDBRepo.EXPECT().GetWhoFollowsPeer(userUUID.String()).Return(nil, repoErr)
		mockLogger.EXPECT().AddFuncName("GetWhoFollowPeer")
		mockLogger.EXPECT().Error("failed to GetWhoFollowPeer from BD")

		ctx := context.WithValue(ctx, config.KeyLogger, mockLogger)

		s := New(mockDBRepo)
		_, err := s.GetWhoFollowPeer(ctx, &friends_proto.GetWhoFollowPeerIn{Uuid: userUUID.String()})
		assert.Error(t, err, repoErr)
	})
}

func TestServer_SetFriends(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDBRepo := NewMockDBRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

	t.Run("should_ok_with_UUID", func(t *testing.T) {
		peer1 := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		peer2 := uuid.Generate()

		mockDBRepo.EXPECT().SetFriend(peer1, peer2.String()).Return(true, nil)
		mockLogger.EXPECT().AddFuncName("SetFriends")
		ctx := context.WithValue(ctx, config.KeyLogger, mockLogger)

		s := New(mockDBRepo)
		res, err := s.SetFriends(ctx, &friends_proto.SetFriendsIn{Peer: peer2.String()})
		assert.NoError(t, err)
		assert.Equal(t, res, &friends_proto.SetFriendsOut{Success: true})
	})

	t.Run("should_is_friends", func(t *testing.T) {
		peer1 := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		peer2 := uuid.Generate()

		mockDBRepo.EXPECT().SetFriend(peer1, peer2.String()).Return(false, nil)
		mockLogger.EXPECT().AddFuncName("SetFriends")
		mockLogger.EXPECT().Info("user already in friends")
		ctx := context.WithValue(ctx, config.KeyLogger, mockLogger)

		s := New(mockDBRepo)
		_, err := s.SetFriends(ctx, &friends_proto.SetFriendsIn{Peer: peer2.String()})
		assert.NoError(t, err)
	})

	t.Run("should_no_UUID", func(t *testing.T) {
		peer1 := ""
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		peer2 := uuid.Generate()

		mockLogger.EXPECT().AddFuncName("SetFriends")
		mockLogger.EXPECT().Error("failed to not found UUID in context")
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		s := New(mockDBRepo)
		_, err := s.SetFriends(ctx, &friends_proto.SetFriendsIn{Peer: peer2.String()})
		assert.Error(t, err)
		assert.EqualError(t, err, "uuid not found in context")
	})

	t.Run("should_repo_err", func(t *testing.T) {
		peer1 := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		peer2 := uuid.Generate()
		repoErr := errors.New("test")

		mockDBRepo.EXPECT().SetFriend(peer1, peer2.String()).Return(false, repoErr)
		mockLogger.EXPECT().AddFuncName("SetFriends")
		mockLogger.EXPECT().Error("failed to SetFriend from BD")
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		s := New(mockDBRepo)
		_, err := s.SetFriends(ctx, &friends_proto.SetFriendsIn{Peer: peer2.String()})
		assert.Error(t, err, repoErr)
	})
}

func TestServer_RemoveFriends(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDBRepo := NewMockDBRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

	t.Run("should_ok_with_UUID", func(t *testing.T) {
		peer1 := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		peer2 := uuid.Generate()

		mockDBRepo.EXPECT().RemoveFriends(peer1, peer2.String()).Return(true, nil)
		mockLogger.EXPECT().AddFuncName("RemoveFriends")
		ctx := context.WithValue(ctx, config.KeyLogger, mockLogger)

		s := New(mockDBRepo)
		res, err := s.RemoveFriends(ctx, &friends_proto.RemoveFriendsIn{Peer: peer2.String()})
		assert.NoError(t, err)
		assert.Equal(t, res, &friends_proto.RemoveFriendsOut{Success: true})
	})

	t.Run("should_no_UUID", func(t *testing.T) {
		peer1 := ""
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		peer2 := uuid.Generate()

		mockLogger.EXPECT().AddFuncName("RemoveFriends")
		mockLogger.EXPECT().Error("failed to not found UUID in context")
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		s := New(mockDBRepo)
		_, err := s.RemoveFriends(ctx, &friends_proto.RemoveFriendsIn{Peer: peer2.String()})
		assert.Error(t, err)
		assert.EqualError(t, err, "uuid not found in context")
	})

	t.Run("should_no_ok_with_UUID", func(t *testing.T) {
		peer1 := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		peer2 := uuid.Generate()

		mockDBRepo.EXPECT().RemoveFriends(peer1, peer2.String()).Return(false, nil)
		mockLogger.EXPECT().AddFuncName("RemoveFriends")
		ctx := context.WithValue(ctx, config.KeyLogger, mockLogger)

		s := New(mockDBRepo)
		_, err := s.RemoveFriends(ctx, &friends_proto.RemoveFriendsIn{Peer: peer2.String()})
		assert.NoError(t, err)
	})

	t.Run("should_repo_err", func(t *testing.T) {
		peer1 := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		peer2 := uuid.Generate()
		repoErr := errors.New("test")

		mockDBRepo.EXPECT().RemoveFriends(peer1, peer2.String()).Return(false, repoErr)
		mockLogger.EXPECT().AddFuncName("RemoveFriends")
		mockLogger.EXPECT().Error("failed to RemoveFriend from BD")
		ctx := context.WithValue(ctx, config.KeyLogger, mockLogger)

		s := New(mockDBRepo)
		_, err := s.RemoveFriends(ctx, &friends_proto.RemoveFriendsIn{Peer: peer2.String()})
		assert.Error(t, err, repoErr)
	})
}

func TestServer_RemoveSubscribe(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDBRepo := NewMockDBRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

	t.Run("should_ok_with_UUID", func(t *testing.T) {
		peer1 := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		peer2 := uuid.Generate()

		mockDBRepo.EXPECT().RemoveSubscribe(peer1, peer2.String()).Return(nil)
		mockLogger.EXPECT().AddFuncName("RemoveSubscribe")
		ctx := context.WithValue(ctx, config.KeyLogger, mockLogger)

		s := New(mockDBRepo)
		res, err := s.RemoveSubscribe(ctx, &friends_proto.RemoveSubscribeIn{Peer: peer2.String()})
		assert.NoError(t, err)
		assert.Equal(t, res, &friends_proto.RemoveSubscribeOut{})
	})

	t.Run("should_no_ok_with_UUID", func(t *testing.T) {
		peer1 := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		peer2 := uuid.Generate()
		repoErr := errors.New("test")

		mockDBRepo.EXPECT().RemoveSubscribe(peer1, peer2.String()).Return(repoErr)
		mockLogger.EXPECT().AddFuncName("RemoveSubscribe")
		mockLogger.EXPECT().Error("failed to RemoveSubscribe from BD")
		ctx := context.WithValue(ctx, config.KeyLogger, mockLogger)

		s := New(mockDBRepo)
		_, err := s.RemoveSubscribe(ctx, &friends_proto.RemoveSubscribeIn{Peer: peer2.String()})
		assert.Error(t, err, repoErr)
	})

	t.Run("should_no_UUID", func(t *testing.T) {
		peer1 := ""
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		peer2 := uuid.Generate()

		mockLogger.EXPECT().AddFuncName("RemoveSubscribe")
		mockLogger.EXPECT().Error("failed to not found UUID in context")
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		s := New(mockDBRepo)
		_, err := s.RemoveSubscribe(ctx, &friends_proto.RemoveSubscribeIn{Peer: peer2.String()})
		assert.Error(t, err)
		assert.EqualError(t, err, "uuid not found in context")
	})
}

func TestServer_InvitePeer(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDBRepo := NewMockDBRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

	t.Run("should_ok_with_uuid_and_email", func(t *testing.T) {
		peerUUID := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peerUUID)
		email := "test@test.ru"

		mockDBRepo.EXPECT().SetInvitePeer(peerUUID, email).Return(nil)
		mockLogger.EXPECT().AddFuncName("SetInvitePeer")
		ctx := context.WithValue(ctx, config.KeyLogger, mockLogger)

		s := New(mockDBRepo)
		_, err := s.SetInvitePeer(ctx, &friends_proto.SetInvitePeerIn{Email: email})
		assert.NoError(t, err)
	})

	t.Run("should_no_ok_with_email", func(t *testing.T) {
		peerUUID := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peerUUID)
		email := ""
		repoErr := errors.New("test")

		mockDBRepo.EXPECT().SetInvitePeer(peerUUID, email).Return(repoErr)
		mockLogger.EXPECT().AddFuncName("SetInvitePeer")
		mockLogger.EXPECT().Error("failed to SetInvitePeer from BD")
		ctx := context.WithValue(ctx, config.KeyLogger, mockLogger)

		s := New(mockDBRepo)
		_, err := s.SetInvitePeer(ctx, &friends_proto.SetInvitePeerIn{Email: email})
		assert.Error(t, err, repoErr)
	})

	t.Run("shold_no_or_with_uuid", func(t *testing.T) {
		peerUUID := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peerUUID)
		email := "test@test.ru"
		repoErr := errors.New("test")

		mockDBRepo.EXPECT().SetInvitePeer(peerUUID, email).Return(repoErr)
		mockLogger.EXPECT().AddFuncName("SetInvitePeer")
		mockLogger.EXPECT().Error("failed to SetInvitePeer from BD")
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		s := New(mockDBRepo)
		_, err := s.SetInvitePeer(ctx, &friends_proto.SetInvitePeerIn{Email: email})
		assert.Error(t, err, repoErr)
	})

	t.Run("should_no_UUID", func(t *testing.T) {
		peer1 := ""
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		email := "test@test.ru"

		mockLogger.EXPECT().AddFuncName("SetInvitePeer")
		mockLogger.EXPECT().Error("failed to not found UUID in context")
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		s := New(mockDBRepo)
		_, err := s.SetInvitePeer(ctx, &friends_proto.SetInvitePeerIn{Email: email})
		assert.Error(t, err)
		assert.EqualError(t, err, "uuid not found in context")
	})
}

func TestServer_GetCountFriends(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

	t.Run("should_ok_with_uuid", func(t *testing.T) {
		ctx := context.Background()
		peerUUID := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peerUUID)
		var subscription int64 = 0
		var subscribers int64 = 0

		mockDBRepo := NewMockDBRepo(ctrl)

		mockDBRepo.EXPECT().GetCountFriends(peerUUID).Return(subscription, subscribers, nil)
		mockLogger.EXPECT().AddFuncName("GetCountFriends")
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		s := New(mockDBRepo)

		res, err := s.GetCountFriends(ctx, &friends_proto.EmptyFriends{})
		assert.NoError(t, err)
		assert.Equal(t, subscription, res.Subscription)
		assert.Equal(t, subscribers, res.Subscribers)
	})

	t.Run("should_no_uuid", func(t *testing.T) {
		ctx := context.Background()
		mockDBRepo := NewMockDBRepo(ctrl)

		mockDBRepo.EXPECT().GetCountFriends(gomock.Any()).Times(0)
		mockLogger.EXPECT().AddFuncName("GetCountFriends")
		mockLogger.EXPECT().Error("failed to not found UUID in context")
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		s := New(mockDBRepo)

		_, err := s.GetCountFriends(ctx, &friends_proto.EmptyFriends{})

		assert.Error(t, err)
		assert.Equal(t, "uuid not found in context", err.Error())
	})

	t.Run("should_error_bd", func(t *testing.T) {
		ctx := context.Background()
		peerUUID := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peerUUID)
		repoErr := errors.New("test")

		mockDBRepo := NewMockDBRepo(ctrl)
		mockDBRepo.EXPECT().GetCountFriends(peerUUID).Return(int64(0), int64(0), repoErr)
		mockLogger.EXPECT().AddFuncName("GetCountFriends")
		mockLogger.EXPECT().Error("failed to GetCountFriends from BD")
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		s := New(mockDBRepo)
		_, err := s.GetCountFriends(ctx, &friends_proto.EmptyFriends{})

		assert.Error(t, err)
		assert.Equal(t, err, repoErr)
	})
}

func TestServer_IsFriendExist(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)
	mockDBRepo := NewMockDBRepo(ctrl)

	t.Run("should_ok_with_uuid", func(t *testing.T) {
		ctx := context.Background()
		userUUID := uuid.Generate().String()
		peerUUID := uuid.Generate().String()

		ctx = context.WithValue(ctx, config.KeyUUID, userUUID)
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		mockLogger.EXPECT().AddFuncName("IsFriendExist")
		mockDBRepo.EXPECT().IsRowFriendExist(userUUID, peerUUID).Return(true, nil)

		s := New(mockDBRepo)

		res, err := s.IsFriendExist(ctx, &friends_proto.IsFriendExistIn{Peer: peerUUID})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.True(t, res.Success)
	})

	t.Run("should_no_ok_UUID_in_context", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)
		peerUUID := uuid.Generate().String()

		mockLogger.EXPECT().AddFuncName("IsFriendExist")
		mockLogger.EXPECT().Error("failed to not found UUID in context")

		s := New(mockDBRepo)

		res, err := s.IsFriendExist(ctx, &friends_proto.IsFriendExistIn{Peer: peerUUID})

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "uuid not found in context")
	})

	t.Run("should_error_bd", func(t *testing.T) {
		ctx := context.Background()
		userUUID := uuid.Generate().String()
		peerUUID := uuid.Generate().String()
		errBD := errors.New("test")

		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)
		ctx = context.WithValue(ctx, config.KeyUUID, userUUID)

		mockLogger.EXPECT().AddFuncName("IsFriendExist")
		mockLogger.EXPECT().Error("failed to IsFriendExist from BD")
		mockDBRepo.EXPECT().IsRowFriendExist(userUUID, peerUUID).Return(false, errBD)

		s := New(mockDBRepo)
		_, err := s.IsFriendExist(ctx, &friends_proto.IsFriendExistIn{Peer: peerUUID})
		assert.Error(t, err)
	})
}
