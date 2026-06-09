package service_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/tmythicator/ticker-rush/backend/internal/apperrors"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/exchange/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/ladder/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/portfolio/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/user/v1"
	app_redis "github.com/tmythicator/ticker-rush/backend/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/backend/internal/service"
	"github.com/tmythicator/ticker-rush/backend/internal/service/mocks"
)

func TestTradeService_BuyStock_Success(t *testing.T) {
	const (
		symbol          string  = "AAPL"
		price           float64 = 100.0
		startBalance    float64 = 1000.0
		userID          int64   = 1
		quantity        float64 = 5.0
		expectedBalance float64 = startBalance - (quantity * price)
	)

	// 1. Setup Market
	mr, _ := miniredis.Run()
	defer mr.Close()

	valkeyClient := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	marketRepo := app_redis.NewMarketRepository(valkeyClient)
	quote := &exchange.Quote{Symbol: symbol, Price: price, Timestamp: time.Now().Unix()}
	bytes, _ := json.Marshal(quote)
	valkeyClient.Set(context.Background(), "market:"+symbol, bytes, 0)

	// 2. Setup Mocks
	mockUserRepo := new(mocks.MockUserRepository)
	mockPortRepo := new(mocks.MockPortfolioRepository)
	mockLadderRepo := new(mocks.MockLadderRepository)
	mockTransactor := new(mocks.MockTransactor)
	mockTx := new(mocks.MockTransaction)

	// 3. Setup Expectations
	ctx := context.Background()

	mockTransactor.On("Begin", mock.Anything).Return(mockTx, nil)
	mockUserRepo.On("WithTx", mockTx).Return(mockUserRepo)
	mockPortRepo.On("WithTx", mockTx).Return(mockPortRepo)
	mockLadderRepo.On("GetActiveLadder", mock.Anything).Return(int64(1), nil)
	mockLadderRepo.On("GetLadder", mock.Anything, int64(1)).Return(&ladder.Ladder{
		Id:             1,
		IsActive:       true,
		StartTime:      timestamppb.New(time.Now().Add(-1 * time.Hour)),
		EndTime:        timestamppb.New(time.Now().Add(1 * time.Hour)),
		InitialBalance: startBalance,
	}, nil)
	mockLadderRepo.On("IsUserInLadder", mock.Anything, int64(1), userID).Return(true, nil)

	initialUser := &user.User{Id: userID, Balance: startBalance}
	mockUserRepo.On("GetUserForUpdate", mock.Anything, userID).Return(initialUser, nil)
	mockUserRepo.On("GetUserBalance", mock.Anything, userID, int64(1)).Return(decimal.NewFromFloat(startBalance), nil)
	mockPortRepo.On("GetPortfolioItemForUpdate", mock.Anything, userID, int64(1), symbol).
		Return(&portfolio.PortfolioItem{StockSymbol: symbol, Quantity: 0}, nil)
	mockUserRepo.On("UpdateUserBalance", mock.Anything, userID, int64(1), mock.MatchedBy(func(d decimal.Decimal) bool {
		v, _ := d.Float64()
		return v == expectedBalance
	})).Return(nil)
	mockPortRepo.On("SetPortfolioItem", mock.Anything, userID, int64(1), symbol, mock.MatchedBy(func(q decimal.Decimal) bool {
		v, _ := q.Float64()
		return v == quantity
	}), mock.MatchedBy(func(p decimal.Decimal) bool {
		v, _ := p.Float64()
		return v == price
	})).Return(nil)
	mockTx.On("Commit", mock.Anything).Return(nil)
	mockTx.On("Rollback", mock.Anything).Return(nil)

	// 4. Execute
	tradeService := service.NewTradeService(mockUserRepo, mockPortRepo, marketRepo, mockLadderRepo, mockTransactor)
	user, err := tradeService.BuyStock(ctx, userID, symbol, quantity)

	// 5. Verify
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedBalance, user.GetBalance())

	mockUserRepo.AssertExpectations(t)
	mockPortRepo.AssertExpectations(t)
	mockTransactor.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestTradeService_BuyStock_InsufficientFunds(t *testing.T) {
	const (
		userID       int64   = 2
		symbol       string  = "ExpensiveStock"
		price        float64 = 10000.0
		quantity     float64 = 1.0
		startBalance float64 = 100.0
	)

	mr, _ := miniredis.Run()
	defer mr.Close()

	rClient := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	marketRepo := app_redis.NewMarketRepository(rClient)
	quote := &exchange.Quote{Symbol: symbol, Price: price, Timestamp: time.Now().Unix()}
	bytes, _ := json.Marshal(quote)
	rClient.Set(context.Background(), "market:"+symbol, bytes, 0)

	mockUserRepo := new(mocks.MockUserRepository)
	mockPortRepo := new(mocks.MockPortfolioRepository)
	mockLadderRepo := new(mocks.MockLadderRepository)
	mockTransactor := new(mocks.MockTransactor)
	mockTx := new(mocks.MockTransaction)

	ctx := context.Background()

	mockTransactor.On("Begin", mock.Anything).Return(mockTx, nil)
	mockTx.On("Rollback", mock.Anything).Return(nil)
	mockUserRepo.On("WithTx", mockTx).Return(mockUserRepo)
	mockPortRepo.On("WithTx", mockTx).Return(mockPortRepo)
	mockLadderRepo.On("GetActiveLadder", mock.Anything).Return(int64(1), nil)
	mockLadderRepo.On("GetLadder", mock.Anything, int64(1)).Return(&ladder.Ladder{
		Id:             1,
		IsActive:       true,
		StartTime:      timestamppb.New(time.Now().Add(-1 * time.Hour)),
		EndTime:        timestamppb.New(time.Now().Add(1 * time.Hour)),
		InitialBalance: startBalance,
	}, nil)
	mockLadderRepo.On("IsUserInLadder", mock.Anything, int64(1), userID).Return(true, nil)

	initialUser := &user.User{Id: userID, Balance: startBalance}
	mockUserRepo.On("GetUserForUpdate", mock.Anything, userID).Return(initialUser, nil)
	mockUserRepo.On("GetUserBalance", mock.Anything, userID, int64(1)).Return(decimal.NewFromFloat(startBalance), nil)

	tradeService := service.NewTradeService(mockUserRepo, mockPortRepo, marketRepo, mockLadderRepo, mockTransactor)
	_, err := tradeService.BuyStock(ctx, userID, symbol, quantity)

	assert.Error(t, err)
	assert.Equal(t, apperrors.ErrInsufficientFunds, err)

	mockTx.AssertNotCalled(t, "Commit", mock.Anything)
	mockTx.AssertCalled(t, "Rollback", mock.Anything)
}

func TestTradeService_BuyStock_MarketClosed(t *testing.T) {
	const (
		userID   int64   = 3
		symbol   string  = "CLOSED_STOCK"
		price    float64 = 150.0
		quantity float64 = 1.0
	)

	// 1. Setup Market with STALE data (1 hour old)
	mr, _ := miniredis.Run()
	defer mr.Close()

	rClient := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	marketRepo := app_redis.NewMarketRepository(rClient)

	// Stale timestamp (1 hour ago)
	staleTime := time.Now().Add(-1 * time.Hour).Unix()
	quote := &exchange.Quote{Symbol: symbol, Price: price, Timestamp: staleTime}
	bytes, _ := json.Marshal(quote)
	rClient.Set(context.Background(), "market:"+symbol, bytes, 0)

	// 2. Setup Mocks (Simulating repositories)
	mockUserRepo := new(mocks.MockUserRepository)
	mockPortRepo := new(mocks.MockPortfolioRepository)
	mockLadderRepo := new(mocks.MockLadderRepository)
	mockTransactor := new(mocks.MockTransactor)

	ctx := context.Background()

	// 3. Execute
	tradeService := service.NewTradeService(mockUserRepo, mockPortRepo, marketRepo, mockLadderRepo, mockTransactor)
	_, err := tradeService.BuyStock(ctx, userID, symbol, quantity)

	// 4. Verify
	assert.Error(t, err)
	assert.Equal(t, apperrors.ErrMarketClosed, err)

	// verify that transaction was NOT started
	mockTransactor.AssertNotCalled(t, "Begin", mock.Anything)
}

func TestTradeService_BuyStock_NotJoined(t *testing.T) {
	const (
		userID   int64   = 4
		symbol   string  = "AAPL"
		price    float64 = 150.0
		quantity float64 = 1.0
	)

	mr, _ := miniredis.Run()
	defer mr.Close()

	rClient := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	marketRepo := app_redis.NewMarketRepository(rClient)
	quote := &exchange.Quote{Symbol: symbol, Price: price, Timestamp: time.Now().Unix()}
	bytes, _ := json.Marshal(quote)
	rClient.Set(context.Background(), "market:"+symbol, bytes, 0)

	mockUserRepo := new(mocks.MockUserRepository)
	mockPortRepo := new(mocks.MockPortfolioRepository)
	mockLadderRepo := new(mocks.MockLadderRepository)
	mockTransactor := new(mocks.MockTransactor)
	mockTx := new(mocks.MockTransaction)

	ctx := context.Background()

	mockTransactor.On("Begin", mock.Anything).Return(mockTx, nil)
	mockTx.On("Rollback", mock.Anything).Return(nil)
	mockUserRepo.On("WithTx", mockTx).Return(mockUserRepo)
	mockPortRepo.On("WithTx", mockTx).Return(mockPortRepo)
	mockUserRepo.On("GetUserForUpdate", mock.Anything, userID).Return(&user.User{Id: userID, Balance: 1000}, nil)
	mockLadderRepo.On("GetActiveLadder", mock.Anything).Return(int64(1), nil)
	mockLadderRepo.On("GetLadder", mock.Anything, int64(1)).Return(&ladder.Ladder{
		Id:             1,
		IsActive:       true,
		StartTime:      timestamppb.New(time.Now().Add(-1 * time.Hour)),
		EndTime:        timestamppb.New(time.Now().Add(1 * time.Hour)),
		InitialBalance: 1000,
	}, nil)
	mockLadderRepo.On("IsUserInLadder", mock.Anything, int64(1), userID).Return(false, nil)

	tradeService := service.NewTradeService(mockUserRepo, mockPortRepo, marketRepo, mockLadderRepo, mockTransactor)
	_, err := tradeService.BuyStock(ctx, userID, symbol, quantity)

	assert.Error(t, err)
	assert.Equal(t, apperrors.ErrNotJoinedLadder, err)
}

func TestTradeService_SellStock_NotJoined(t *testing.T) {
	const (
		userID   int64   = 4
		symbol   string  = "AAPL"
		price    float64 = 150.0
		quantity float64 = 1.0
	)

	mr, _ := miniredis.Run()
	defer mr.Close()

	rClient := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	marketRepo := app_redis.NewMarketRepository(rClient)
	quote := &exchange.Quote{Symbol: symbol, Price: price, Timestamp: time.Now().Unix()}
	bytes, _ := json.Marshal(quote)
	rClient.Set(context.Background(), "market:"+symbol, bytes, 0)

	mockUserRepo := new(mocks.MockUserRepository)
	mockPortRepo := new(mocks.MockPortfolioRepository)
	mockLadderRepo := new(mocks.MockLadderRepository)
	mockTransactor := new(mocks.MockTransactor)
	mockTx := new(mocks.MockTransaction)

	ctx := context.Background()

	mockTransactor.On("Begin", mock.Anything).Return(mockTx, nil)
	mockTx.On("Rollback", mock.Anything).Return(nil)
	mockUserRepo.On("WithTx", mockTx).Return(mockUserRepo)
	mockPortRepo.On("WithTx", mockTx).Return(mockPortRepo)
	mockLadderRepo.On("GetActiveLadder", mock.Anything).Return(int64(1), nil)
	mockLadderRepo.On("GetLadder", mock.Anything, int64(1)).Return(&ladder.Ladder{
		Id:             1,
		IsActive:       true,
		StartTime:      timestamppb.New(time.Now().Add(-1 * time.Hour)),
		EndTime:        timestamppb.New(time.Now().Add(1 * time.Hour)),
		InitialBalance: 1000,
	}, nil)
	mockLadderRepo.On("IsUserInLadder", mock.Anything, int64(1), userID).Return(false, nil)

	tradeService := service.NewTradeService(mockUserRepo, mockPortRepo, marketRepo, mockLadderRepo, mockTransactor)
	_, err := tradeService.SellStock(ctx, userID, symbol, quantity)

	assert.Error(t, err)
	assert.Equal(t, apperrors.ErrNotJoinedLadder, err)
}

func TestTradeService_BuyStock_LadderNotActive(t *testing.T) {
	const (
		userID   int64   = 5
		symbol   string  = "AAPL"
		price    float64 = 150.0
		quantity float64 = 1.0
	)

	mr, _ := miniredis.Run()
	defer mr.Close()

	rClient := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	marketRepo := app_redis.NewMarketRepository(rClient)
	quote := &exchange.Quote{Symbol: symbol, Price: price, Timestamp: time.Now().Unix()}
	bytes, _ := json.Marshal(quote)
	rClient.Set(context.Background(), "market:"+symbol, bytes, 0)

	mockUserRepo := new(mocks.MockUserRepository)
	mockPortRepo := new(mocks.MockPortfolioRepository)
	mockLadderRepo := new(mocks.MockLadderRepository)
	mockTransactor := new(mocks.MockTransactor)

	ctx := context.Background()

	mockLadderRepo.On("GetActiveLadder", mock.Anything).Return(int64(1), nil)
	// Mock an inactive ladder (ends in the past)
	mockLadderRepo.On("GetLadder", mock.Anything, int64(1)).Return(&ladder.Ladder{
		Id:             1,
		IsActive:       true,
		StartTime:      timestamppb.New(time.Now().Add(-2 * time.Hour)),
		EndTime:        timestamppb.New(time.Now().Add(-1 * time.Hour)),
		InitialBalance: 1000,
	}, nil)

	tradeService := service.NewTradeService(mockUserRepo, mockPortRepo, marketRepo, mockLadderRepo, mockTransactor)
	_, err := tradeService.BuyStock(ctx, userID, symbol, quantity)

	assert.Error(t, err)
	assert.Equal(t, apperrors.ErrLadderNotActive, err)
}

