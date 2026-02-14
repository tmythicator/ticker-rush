package service_test

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/tmythicator/ticker-rush/server/internal/proto/exchange/v1"
	"github.com/tmythicator/ticker-rush/server/internal/proto/leaderboard/v1"
	"github.com/tmythicator/ticker-rush/server/internal/proto/user/v1"
	"github.com/tmythicator/ticker-rush/server/internal/service"
	"github.com/tmythicator/ticker-rush/server/internal/service/mocks"
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

	lbService := service.NewLeaderBoardService(mockUserRepo, mockPortfolioRepo, mockMarketRepo, redisClient)

	ctx := context.Background()

	users := []*user.User{
		{Id: 1, Balance: 1000, FirstName: "Alice"},
		{Id: 2, Balance: 2000, FirstName: "Bob"},
	}

	portfolio1 := []*user.PortfolioItem{
		{StockSymbol: "AAPL", Quantity: 10},
	}
	portfolio2 := []*user.PortfolioItem{
		{StockSymbol: "GOOG", Quantity: 5},
	}

	quoteAAPL := &exchange.Quote{Symbol: "AAPL", Price: 150.0}
	quoteGOOG := &exchange.Quote{Symbol: "GOOG", Price: 200.0}

	mockUserRepo.On("GetUsers", ctx).Return(users, nil)

	mockPortfolioRepo.On("GetPortfolio", ctx, int64(1)).Return(portfolio1, nil)
	mockPortfolioRepo.On("GetPortfolio", ctx, int64(2)).Return(portfolio2, nil)

	mockMarketRepo.On("GetQuote", ctx, "AAPL").Return(quoteAAPL, nil)
	mockMarketRepo.On("GetQuote", ctx, "GOOG").Return(quoteGOOG, nil)

	err = lbService.UpdateLeaderboard(ctx)
	assert.NoError(t, err)

	// Alice: 1000 + (10 * 150) = 1000 + 1500 = 2500
	// Bob:   2000 + (5 * 200) = 2000 + 1000 = 3000

	scoreAlice, err := redisClient.ZScore(ctx, "leaderboard", "1").Result()
	assert.NoError(t, err)
	assert.Equal(t, 2500.0, scoreAlice)

	scoreBob, err := redisClient.ZScore(ctx, "leaderboard", "2").Result()
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
	lbService := service.NewLeaderBoardService(mockUserRepo, nil, nil, redisClient)

	ctx := context.Background()

	// Setup Redis data
	redisClient.ZAdd(ctx, "leaderboard", redis.Z{Score: 3000, Member: "2"})
	redisClient.ZAdd(ctx, "leaderboard", redis.Z{Score: 2500, Member: "1"})
	redisClient.ZAdd(ctx, "leaderboard", redis.Z{Score: 1500, Member: "3"})

	mockUserRepo.On("GetUser", ctx, int64(2)).Return(&user.User{Id: 2, FirstName: "Bob", LastName: "B"}, nil)
	mockUserRepo.On("GetUser", ctx, int64(1)).Return(&user.User{Id: 1, FirstName: "Alice", LastName: "A"}, nil)
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

	assert.Equal(t, int64(2), resp.Entries[0].UserId)
	assert.Equal(t, int32(1), resp.Entries[0].Rank)
	assert.Equal(t, 3000.0, resp.Entries[0].TotalNetWorth)
	assert.Equal(t, "Bob", resp.Entries[0].FirstName)

	assert.Equal(t, int64(1), resp.Entries[1].UserId)
	assert.Equal(t, int32(2), resp.Entries[1].Rank)
	assert.Equal(t, 2500.0, resp.Entries[1].TotalNetWorth)
	assert.Equal(t, "Alice", resp.Entries[1].FirstName)

	// Verify User 3 was removed from Redis
	exists := redisClient.ZScore(ctx, "leaderboard", "3").Val()
	assert.Equal(t, 0.0, exists)
}
