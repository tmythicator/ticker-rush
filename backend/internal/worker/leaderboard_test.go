package worker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tmythicator/ticker-rush/backend/internal/domain"
	"github.com/tmythicator/ticker-rush/backend/internal/service"
	"github.com/tmythicator/ticker-rush/backend/internal/service/mocks"
)

func TestLeaderboardWorker_Start(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockPortRepo := new(mocks.MockPortfolioRepository)
	mockMarketRepo := new(mocks.MockMarketRepository)
	mockLeaderboardRepo := new(mocks.MockLeaderboardRepository)
	mockLadderRepo := new(mocks.MockLadderRepository)

	// Mock minimal calls for UpdateLeaderboard
	// Return active ladder ID 1 and empty users slice so UpdateLeaderboard completes quickly with no-op
	mockLadderRepo.On("GetActiveLadder", mock.Anything).Return(int64(1), nil)
	mockUserRepo.On("GetUsers", mock.Anything).Return([]*domain.User{}, nil)
	mockLeaderboardRepo.On("SetLastUpdate", mock.Anything, int64(1), mock.Anything).Return(nil)

	lbService := service.NewLeaderBoardService(
		mockUserRepo,
		mockPortRepo,
		mockMarketRepo,
		mockLadderRepo,
		mockLeaderboardRepo,
	)

	// Run with 10ms interval, cancel context after 25ms to allow initial update + 2 ticks
	w := NewLeaderboardWorker(lbService, 10*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- w.Start(ctx)
	}()

	// Wait 25ms and then cancel context to stop the worker
	time.Sleep(25 * time.Millisecond)
	cancel()

	err := <-errChan
	assert.ErrorIs(t, err, context.Canceled)

	// Verify that UpdateLeaderboard was called at least once
	mockLadderRepo.AssertCalled(t, "GetActiveLadder", mock.Anything)
	mockUserRepo.AssertCalled(t, "GetUsers", mock.Anything)
}
