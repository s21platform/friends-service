package mocks

import (
	"context"
	"errors"
	"github.com/docker/distribution/uuid"
	"github.com/golang/mock/gomock"
	"github.com/s21platform/friends-proto/friends-proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServer_GetPeerFollow(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDbRepo := NewMockDbRepo(ctrl)

	t.Run("should_ok_with_UUID", func(t *testing.T) {
		userUUID := uuid.Generate()
		followersUUID := []string{
			uuid.Generate().String(),
			uuid.Generate().String(),
		}
		mockDbRepo.EXPECT().GetPeerFollows(userUUID.String()).Return(followersUUID, nil)

		s := New(mockDbRepo)
		res, err := s.GetPeerFollow(ctx, &friends_proto.GetPeerFollowIn{Uuid: userUUID.String()})
		assert.NoError(t, err)
		assert.Equal(t, res, &friends_proto.GetPeerFollowOut{Subscription: []*friends_proto.Peer{
			&friends_proto.Peer{Uuid: followersUUID[0]},
			&friends_proto.Peer{Uuid: followersUUID[1]},
		}})

	})

	t.Run("should_repo_err", func(t *testing.T) {
		userUUID := uuid.Generate()
		repoErr := errors.New("test")

		mockDbRepo.EXPECT().GetPeerFollows(userUUID.String()).Return(nil, repoErr)

		s := New(mockDbRepo)
		_, err := s.GetPeerFollow(ctx, &friends_proto.GetPeerFollowIn{Uuid: userUUID.String()})
		assert.Error(t, err, repoErr)

	})
}
