package service

import (
	"context"
	"errors"
	"testing"

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

	t.Run("should_ok_with_UUID", func(t *testing.T) {
		userUUID := uuid.Generate()
		followersUUID := []string{
			uuid.Generate().String(),
			uuid.Generate().String(),
		}
		mockDBRepo.EXPECT().GetPeerFollows(userUUID.String()).Return(followersUUID, nil)

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

	t.Run("should_ok_with_UUID", func(t *testing.T) {
		userUUID := uuid.Generate()
		followersUUID := []string{
			uuid.Generate().String(),
			uuid.Generate().String(),
		}
		mockDBRepo.EXPECT().GetWhoFollowsPeer(userUUID.String()).Return(followersUUID, nil)

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

	t.Run("should_ok_with_UUID", func(t *testing.T) {
		peer1 := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		peer2 := uuid.Generate()

		mockDBRepo.EXPECT().SetFriend(peer1, peer2.String()).Return(true, nil)

		s := New(mockDBRepo)
		res, err := s.SetFriends(ctx, &friends_proto.SetFriendsIn{Peer: peer2.String()})
		assert.NoError(t, err)
		assert.Equal(t, res, &friends_proto.SetFriendsOut{Success: true})
	})

	t.Run("should_no_ok_with_UUID", func(t *testing.T) {
		peer1 := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		peer2 := uuid.Generate()

		mockDBRepo.EXPECT().SetFriend(peer1, peer2.String()).Return(false, nil)

		s := New(mockDBRepo)
		_, err := s.SetFriends(ctx, &friends_proto.SetFriendsIn{Peer: peer2.String()})
		assert.NoError(t, err)
	})

	t.Run("should_repo_err", func(t *testing.T) {
		peer1 := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		peer2 := uuid.Generate()
		repoErr := errors.New("test")

		mockDBRepo.EXPECT().SetFriend(peer1, peer2.String()).Return(false, repoErr)

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

	t.Run("should_ok_with_UUID", func(t *testing.T) {
		peer1 := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		peer2 := uuid.Generate()

		mockDBRepo.EXPECT().RemoveFriends(peer1, peer2.String()).Return(true, nil)

		s := New(mockDBRepo)
		res, err := s.RemoveFriends(ctx, &friends_proto.RemoveFriendsIn{Peer: peer2.String()})
		assert.NoError(t, err)
		assert.Equal(t, res, &friends_proto.RemoveFriendsOut{Success: true})
	})

	t.Run("should_no_ok_with_UUID", func(t *testing.T) {
		peer1 := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		peer2 := uuid.Generate()

		mockDBRepo.EXPECT().RemoveFriends(peer1, peer2.String()).Return(false, nil)

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

	t.Run("should_ok_with_UUID", func(t *testing.T) {
		peer1 := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peer1)
		peer2 := uuid.Generate()

		mockDBRepo.EXPECT().RemoveSubscribe(peer1, peer2.String()).Return(nil)
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

		s := New(mockDBRepo)
		_, err := s.RemoveSubscribe(ctx, &friends_proto.RemoveSubscribeIn{Peer: peer2.String()})
		assert.Error(t, err, repoErr)
	})
}

func TestServer_InvitePeer(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDBRepo := NewMockDBRepo(ctrl)

	t.Run("should_ok_with_uuid_and_email", func(t *testing.T) {
		peerUUID := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peerUUID)
		email := "test@test.ru"

		mockDBRepo.EXPECT().SetInvitePeer(peerUUID, email).Return(nil)
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
		s := New(mockDBRepo)
		_, err := s.SetInvitePeer(ctx, &friends_proto.SetInvitePeerIn{Email: email})
		assert.Error(t, err, repoErr)
	})
}

func TestServer_GetCountFriends(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDBRepo := NewMockDBRepo(ctrl)

	t.Run("should_ok_with_uuid", func(t *testing.T) {
		peerUUID := uuid.Generate().String()
		ctx = context.WithValue(ctx, config.KeyUUID, peerUUID)
		var subscription int64 = 0
		var subscribers int64 = 0

		mockDBRepo.EXPECT().GetCountFriends(peerUUID).Return(subscription, subscribers, nil)

		s := New(mockDBRepo)

		res, err := s.GetCountFriends(ctx, &friends_proto.EmptyFriends{})
		assert.NoError(t, err)
		assert.Equal(t, subscription, res.Subscription)
		assert.Equal(t, subscribers, res.Subscribers)
	})
	t.Run("should_no_uuid", func(t *testing.T) {
		peerUUID := ""
		ctx = context.WithValue(ctx, config.KeyUUID, peerUUID)
		var subscription int64 = 0
		var subscribers int64 = 0
		repoErr := errors.New("test")

		mockDBRepo.EXPECT().GetCountFriends(peerUUID).Return(subscription, subscribers, repoErr)

		s := New(mockDBRepo)

		_, err := s.GetCountFriends(ctx, &friends_proto.EmptyFriends{})
		assert.Error(t, err, repoErr)
	})
}
