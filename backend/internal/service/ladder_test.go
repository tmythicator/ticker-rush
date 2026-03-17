package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

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
		mockTransactor := new(mocks.MockTransactor)
		mockTx := new(mocks.MockTransaction)

		mockRepo.On("GetActiveLadder", ctx).Return(ladderID, nil)
		mockRepo.On("IsUserInLadder", ctx, ladderID, userID).Return(false, nil)
		mockTransactor.On("Begin", ctx).Return(mockTx, nil)
		mockRepo.On("WithTx", mockTx).Return(mockRepo)
		mockRepo.On("JoinLadder", ctx, ladderID, userID).Return(nil)
		mockTx.On("Commit", ctx).Return(nil)
		mockTx.On("Rollback", ctx).Return(nil)

		s := service.NewLadderService(mockRepo, mockTransactor)
		err := s.JoinLadder(ctx, userID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockTransactor.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("AlreadyJoined", func(t *testing.T) {
		mockRepo := new(mocks.MockLadderRepository)
		mockTransactor := new(mocks.MockTransactor)

		mockRepo.On("GetActiveLadder", ctx).Return(ladderID, nil)
		mockRepo.On("IsUserInLadder", ctx, ladderID, userID).Return(true, nil)

		s := service.NewLadderService(mockRepo, mockTransactor)
		err := s.JoinLadder(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrAlreadyJoinedLadder, err)
		mockTransactor.AssertNotCalled(t, "Begin", mock.Anything)
	})

	t.Run("TransactionFailure", func(t *testing.T) {
		mockRepo := new(mocks.MockLadderRepository)
		mockTransactor := new(mocks.MockTransactor)
		mockTx := new(mocks.MockTransaction)

		mockRepo.On("GetActiveLadder", ctx).Return(ladderID, nil)
		mockRepo.On("IsUserInLadder", ctx, ladderID, userID).Return(false, nil)
		mockTransactor.On("Begin", ctx).Return(mockTx, nil)
		mockRepo.On("WithTx", mockTx).Return(mockRepo)
		mockRepo.On("JoinLadder", ctx, ladderID, userID).Return(assert.AnError)
		mockTx.On("Rollback", ctx).Return(nil)

		s := service.NewLadderService(mockRepo, mockTransactor)
		err := s.JoinLadder(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
		mockTx.AssertNotCalled(t, "Commit", mock.Anything)
	})
}
