package worker_test

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/tmythicator/ticker-rush/backend/internal/gen/sqlc"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/exchange/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/ladder/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/portfolio/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/service/mocks"
	"github.com/tmythicator/ticker-rush/backend/internal/worker"
)

func TestLadderLifecycleWorker_RunOnce_ExpiredLadders(t *testing.T) {
	mockLadderRepo := new(mocks.MockLadderRepository)
	mockPortRepo := new(mocks.MockPortfolioRepository)
	mockMarketRepo := new(mocks.MockMarketRepository)

	ctx := context.Background()
	now := time.Now()

	// 1. Setup expired ladder
	expiredLadder := &ladder.Ladder{
		Id:             10,
		Name:           "Expired Season 1",
		IsActive:       true,
		StartTime:      timestamppb.New(now.Add(-2 * time.Hour)),
		EndTime:        timestamppb.New(now.Add(-1 * time.Hour)),
		InitialBalance: 10000.0,
	}

	mockLadderRepo.On("GetExpiredActiveLadders", mock.Anything, mock.Anything).Return([]*ladder.Ladder{expiredLadder}, nil)

	// 2. Setup participants: User 101 and User 102
	participants := []sqlc.LadderParticipant{
		{
			LadderID: 10,
			UserID:   101,
			Balance:  decimal.NewFromFloat(1000.0), // Liquid cash
		},
		{
			LadderID: 10,
			UserID:   102,
			Balance:  decimal.NewFromFloat(2000.0), // Liquid cash
		},
	}
	mockLadderRepo.On("GetLadderParticipants", mock.Anything, int64(10)).Return(participants, nil)

	// User 101 has 10 AAPL. User 102 has 5 AAPL.
	mockPortRepo.On("GetPortfolio", mock.Anything, int64(101), int64(10)).Return([]*portfolio.PortfolioItem{
		{
			StockSymbol: "AAPL",
			Quantity:    10.0,
		},
	}, nil)

	mockPortRepo.On("GetPortfolio", mock.Anything, int64(102), int64(10)).Return([]*portfolio.PortfolioItem{
		{
			StockSymbol: "AAPL",
			Quantity:    5.0,
		},
	}, nil)

	// AAPL price is $150
	mockMarketRepo.On("GetQuote", mock.Anything, "AAPL").Return(&exchange.Quote{
		Symbol: "AAPL",
		Price:  150.0,
	}, nil)

	// Calculations:
	// User 101: Cash (1000) + 10 * 150 = 2500 -> Rank 1
	// User 102: Cash (2000) + 5 * 150 = 2750 -> Rank 1 (Wait: 2750 > 2500, so User 102 is Rank 1, User 101 is Rank 2)

	// We expect updates to database
	mockLadderRepo.On("InsertLadderParticipant", mock.Anything, int64(10), int64(102), mock.MatchedBy(func(d decimal.Decimal) bool {
		val, _ := d.Float64()

		return val == 2750.0
	}), int32(1)).Return(nil)

	mockLadderRepo.On("InsertLadderParticipant", mock.Anything, int64(10), int64(101), mock.MatchedBy(func(d decimal.Decimal) bool {
		val, _ := d.Float64()

		return val == 2500.0
	}), int32(2)).Return(nil)

	mockLadderRepo.On("DeleteLadderPortfolioItemsByLadder", mock.Anything, int64(10)).Return(nil)
	mockLadderRepo.On("PruneLadderParticipants", mock.Anything, int64(10), int32(20)).Return(nil)
	mockLadderRepo.On("UpdateLadderStatus", mock.Anything, int64(10), false).Return(nil)

	// Setup no pending ladders
	mockLadderRepo.On("GetPendingLaddersToActivate", mock.Anything, mock.Anything).Return([]*ladder.Ladder{}, nil)

	// Create and run worker
	w := worker.NewLadderLifecycleWorker(mockLadderRepo, mockPortRepo, mockMarketRepo, 10*time.Millisecond)
	w.RunOnce(ctx)

	mockLadderRepo.AssertExpectations(t)
	mockPortRepo.AssertExpectations(t)
	mockMarketRepo.AssertExpectations(t)
}

func TestLadderLifecycleWorker_RunOnce_ActivatePendingLadders(t *testing.T) {
	mockLadderRepo := new(mocks.MockLadderRepository)
	mockPortRepo := new(mocks.MockPortfolioRepository)
	mockMarketRepo := new(mocks.MockMarketRepository)

	ctx := context.Background()
	now := time.Now()

	// No expired ladders
	mockLadderRepo.On("GetExpiredActiveLadders", mock.Anything, mock.Anything).Return([]*ladder.Ladder{}, nil)

	// Setup pending ladder
	pendingLadder := &ladder.Ladder{
		Id:             20,
		Name:           "Pending Season 2",
		IsActive:       false,
		StartTime:      timestamppb.New(now.Add(-10 * time.Minute)),
		EndTime:        timestamppb.New(now.Add(1 * time.Hour)),
		InitialBalance: 5000.0,
	}
	mockLadderRepo.On("GetPendingLaddersToActivate", mock.Anything, mock.Anything).Return([]*ladder.Ladder{pendingLadder}, nil)
	mockLadderRepo.On("UpdateLadderStatus", mock.Anything, int64(20), true).Return(nil)

	// Create and run worker
	w := worker.NewLadderLifecycleWorker(mockLadderRepo, mockPortRepo, mockMarketRepo, 10*time.Millisecond)
	w.RunOnce(ctx)

	mockLadderRepo.AssertExpectations(t)
}
