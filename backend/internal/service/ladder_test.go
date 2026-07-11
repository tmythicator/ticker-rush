package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tmythicator/ticker-rush/backend/internal/apperrors"
	"github.com/tmythicator/ticker-rush/backend/internal/service"
	"github.com/tmythicator/ticker-rush/backend/internal/service/mocks"
)

func TestLadderService_JoinLadder(t *testing.T) {
	const (
		ladderID int64 = 1
		userID   int64 = 100
	)

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mocks.MockLadderRepository)

		mockRepo.On("GetActiveLadder", ctx).Return(ladderID, nil)
		mockRepo.On("JoinLadder", ctx, ladderID, userID).Return(nil)

		s := service.NewLadder(mockRepo)
		err := s.JoinLadder(ctx, userID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("AlreadyJoined", func(t *testing.T) {
		mockRepo := new(mocks.MockLadderRepository)

		mockRepo.On("GetActiveLadder", ctx).Return(ladderID, nil)
		mockRepo.On("JoinLadder", ctx, ladderID, userID).Return(apperrors.ErrAlreadyJoinedLadder)

		s := service.NewLadder(mockRepo)
		err := s.JoinLadder(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrAlreadyJoinedLadder, err)
		mockRepo.AssertExpectations(t)
	})
}
