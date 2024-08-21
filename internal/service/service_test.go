package service

import (
	"context"
	"errors"
	"testing"

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

	mockDbRepo := NewMockDBRepo(ctrl)

	t.Run("should_ok_with_UUID", func(t *testing.T) {
		userUUID := uuid.Generate()
		followersUUID := []string{
			uuid.Generate().String(),
			uuid.Generate().String(),
		}
		mockDbRepo.EXPECT().GetWhoFollowsPeer(userUUID.String()).Return(followersUUID, nil)

		s := New(mockDbRepo)
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

		mockDbRepo.EXPECT().GetWhoFollowsPeer(userUUID.String()).Return(nil, repoErr)

		s := New(mockDbRepo)
		_, err := s.GetWhoFollowPeer(ctx, &friends_proto.GetWhoFollowPeerIn{Uuid: userUUID.String()})
		assert.Error(t, err, repoErr)
	})
}

func TestServer_SetFriends(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbRepo := NewMockDBRepo(ctrl)

	t.Run("should_ok_with_UUID", func(t *testing.T) {
		peer1 := uuid.Generate()
		peer2 := uuid.Generate()

		mockDbRepo.EXPECT().SetFriend(peer1.String(), peer2.String()).Return(true, nil)

		s := New(mockDbRepo)
		res, err := s.SetFriends(ctx, &friends_proto.SetFriendsIn{Peer_1: peer1.String(), Peer_2: peer2.String()})
		assert.NoError(t, err)
		assert.Equal(t, res, &friends_proto.SetFriendsOut{Success: true})
	})

	t.Run("should_no_ok_with_UUID", func(t *testing.T) {
		peer1 := uuid.Generate()
		peer2 := uuid.Generate()

		mockDbRepo.EXPECT().SetFriend(peer1.String(), peer2.String()).Return(false, nil)

		s := New(mockDbRepo)
		_, err := s.SetFriends(ctx, &friends_proto.SetFriendsIn{Peer_1: peer1.String(), Peer_2: peer2.String()})
		assert.NoError(t, err)
	})

	t.Run("should_repo_err", func(t *testing.T) {
		peer1 := uuid.Generate()
		peer2 := uuid.Generate()
		repoErr := errors.New("test")

		mockDbRepo.EXPECT().SetFriend(peer1.String(), peer2.String()).Return(false, repoErr)

		s := New(mockDbRepo)
		_, err := s.SetFriends(ctx, &friends_proto.SetFriendsIn{Peer_1: peer1.String(), Peer_2: peer2.String()})
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
		peer1 := uuid.Generate()
		peer2 := uuid.Generate()

		mockDBRepo.EXPECT().RemoveSubscribe(peer1.String(), peer2.String()).Return(nil)
		s := New(mockDBRepo)
		res, err := s.RemoveSubscribe(ctx, &friends_proto.RemoveSubscribeIn{Peer_1: peer1.String(), Peer_2: peer2.String()})
		assert.NoError(t, err)
		assert.Equal(t, res, &friends_proto.RemoveSubscribeOut{})
	})

	t.Run("should_no_ok_with_UUID", func(t *testing.T) {
		peer1 := uuid.Generate()
		peer2 := uuid.Generate()
		repoErr := errors.New("test")

		mockDBRepo.EXPECT().RemoveSubscribe(peer1.String(), peer2.String()).Return(repoErr)

		s := New(mockDBRepo)
		_, err := s.RemoveSubscribe(ctx, &friends_proto.RemoveSubscribeIn{Peer_1: peer1.String(), Peer_2: peer2.String()})
		assert.Error(t, err, repoErr)
	})
}
