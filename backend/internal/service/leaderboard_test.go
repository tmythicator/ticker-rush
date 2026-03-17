package service_test

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/tmythicator/ticker-rush/backend/internal/proto/exchange/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/leaderboard/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/portfolio/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/user/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/service"
	"github.com/tmythicator/ticker-rush/backend/internal/service/mocks"
)

func TestLeaderBoardService_UpdateLeaderboard(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	mockUserRepo := new(mocks.MockUserRepository)
	mockPortfolioRepo := new(mocks.MockPortfolioRepository)
	mockMarketRepo := new(mocks.MockMarketRepository)
	mockLadderRepo := new(mocks.MockLadderRepository)

	lbService := service.NewLeaderBoardService(mockUserRepo, mockPortfolioRepo, mockMarketRepo, mockLadderRepo, redisClient)

	ctx := context.Background()

	users := []*user.User{
		{Id: 1, Balance: 1000, FirstName: "Alice"},
		{Id: 2, Balance: 2000, FirstName: "Bob"},
	}

	portfolio1 := []*portfolio.PortfolioItem{
		{StockSymbol: "AAPL", Quantity: 10},
	}
	portfolio2 := []*portfolio.PortfolioItem{
		{StockSymbol: "GOOG", Quantity: 5},
	}

	quoteAAPL := &exchange.Quote{Symbol: "AAPL", Price: 150.0}
	quoteGOOG := &exchange.Quote{Symbol: "GOOG", Price: 200.0}

	mockLadderRepo.On("GetActiveLadder", ctx).Return(int64(1), nil)
	mockUserRepo.On("GetUsers", ctx).Return(users, nil)

	mockUserRepo.On("GetUserBalance", ctx, int64(1), int64(1)).Return(1000.0, nil)
	mockUserRepo.On("GetUserBalance", ctx, int64(2), int64(1)).Return(2000.0, nil)

	mockPortfolioRepo.On("GetPortfolio", ctx, int64(1), int64(1)).Return(portfolio1, nil)
	mockPortfolioRepo.On("GetPortfolio", ctx, int64(2), int64(1)).Return(portfolio2, nil)

	mockMarketRepo.On("GetQuote", ctx, "AAPL").Return(quoteAAPL, nil)
	mockMarketRepo.On("GetQuote", ctx, "GOOG").Return(quoteGOOG, nil)

	err = lbService.UpdateLeaderboard(ctx)
	assert.NoError(t, err)

	// Alice: 1000 + (10 * 150) = 1000 + 1500 = 2500
	// Bob:   2000 + (5 * 200) = 2000 + 1000 = 3000

	scoreAlice, err := redisClient.ZScore(ctx, "leaderboard:1", "1").Result()
	assert.NoError(t, err)
	assert.Equal(t, 2500.0, scoreAlice)

	scoreBob, err := redisClient.ZScore(ctx, "leaderboard:1", "2").Result()
	assert.NoError(t, err)
	assert.Equal(t, 3000.0, scoreBob)
}

func TestLeaderBoardService_GetLeaderboard(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	mockUserRepo := new(mocks.MockUserRepository)
	mockLadderRepo := new(mocks.MockLadderRepository)
	lbService := service.NewLeaderBoardService(mockUserRepo, nil, nil, mockLadderRepo, redisClient)

	ctx := context.Background()

	mockLadderRepo.On("GetActiveLadder", ctx).Return(int64(1), nil)

	// Setup Redis data
	redisClient.ZAdd(ctx, "leaderboard:1", redis.Z{Score: 3000, Member: "2"})
	redisClient.ZAdd(ctx, "leaderboard:1", redis.Z{Score: 2500, Member: "1"})
	redisClient.ZAdd(ctx, "leaderboard:1", redis.Z{Score: 1500, Member: "3"})

	mockUserRepo.On("GetUser", ctx, int64(2)).Return(&user.User{Id: 2, FirstName: "Bob", LastName: "B", IsPublic: false}, nil)
	mockUserRepo.On("GetUser", ctx, int64(1)).Return(&user.User{
		Id:        1,
		Username:  "user1",
		FirstName: "First1",
		LastName:  "Last1",
		IsPublic:  true,
	}, nil)
	// User 3 is missing (testing cleanup)
	mockUserRepo.On("GetUser", ctx, int64(3)).Return(nil, assert.AnError)

	req := &leaderboard.GetLeaderboardRequest{
		Limit:  10,
		Offset: 0,
	}

	resp, err := lbService.GetLeaderboard(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// User 3 should be removed, so total count should be 2.
	assert.Equal(t, int32(2), resp.TotalCount)
	assert.Len(t, resp.Entries, 2)

	assert.Equal(t, int64(2), resp.Entries[0].User.Id)
	assert.Equal(t, int32(1), resp.Entries[0].Rank)
	assert.Equal(t, 3000.0, resp.Entries[0].Score)
	assert.Equal(t, "Classified", resp.Entries[0].User.Username)
	assert.Equal(t, "", resp.Entries[0].User.FirstName)
	assert.Equal(t, "", resp.Entries[0].User.LastName)
	assert.False(t, resp.Entries[0].User.IsPublic)

	assert.Equal(t, int64(1), resp.Entries[1].User.Id)
	assert.Equal(t, int32(2), resp.Entries[1].Rank)
	assert.Equal(t, 2500.0, resp.Entries[1].Score)
	assert.Equal(t, "First1", resp.Entries[1].User.FirstName)
	assert.Equal(t, "Last1", resp.Entries[1].User.LastName)
	assert.True(t, resp.Entries[1].User.IsPublic)

	// Verify User 3 was removed from Redis
	exists := redisClient.ZScore(ctx, "leaderboard:1", "3").Val()
	assert.Equal(t, 0.0, exists)
}
