package service_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tmythicator/ticker-rush/server/internal/apperrors"
	"github.com/tmythicator/ticker-rush/server/internal/proto/exchange/v1"
	"github.com/tmythicator/ticker-rush/server/internal/proto/user/v1"
	app_redis "github.com/tmythicator/ticker-rush/server/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/server/internal/service"
	"github.com/tmythicator/ticker-rush/server/internal/service/mocks"
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
	mockTransactor := new(mocks.MockTransactor)
	mockTx := new(mocks.MockTransaction)

	// 3. Setup Expectations
	ctx := context.Background()

	mockTransactor.On("Begin", mock.Anything).Return(mockTx, nil)
	mockUserRepo.On("WithTx", mockTx).Return(mockUserRepo)
	mockPortRepo.On("WithTx", mockTx).Return(mockPortRepo)

	initialUser := &user.User{Id: userID, Balance: startBalance}
	mockUserRepo.On("GetUserForUpdate", mock.Anything, userID).Return(initialUser, nil)
	mockPortRepo.On("GetPortfolioItemForUpdate", mock.Anything, userID, symbol).
		Return(&user.PortfolioItem{}, assert.AnError)
	mockUserRepo.On("SaveUser", mock.Anything, mock.MatchedBy(func(u *user.User) bool {
		return u.GetBalance() == expectedBalance
	})).Return(nil)
	mockPortRepo.On("SetPortfolioItem", mock.Anything, userID, symbol, quantity, price).Return(nil)
	mockTx.On("Commit", mock.Anything).Return(nil)
	mockTx.On("Rollback", mock.Anything).Return(nil)

	// 4. Execute
	tradeService := service.NewTradeService(mockUserRepo, mockPortRepo, marketRepo, mockTransactor)
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
	mockTransactor := new(mocks.MockTransactor)
	mockTx := new(mocks.MockTransaction)

	ctx := context.Background()

	mockTransactor.On("Begin", mock.Anything).Return(mockTx, nil)
	mockTx.On("Rollback", mock.Anything).Return(nil)
	mockUserRepo.On("WithTx", mockTx).Return(mockUserRepo)
	mockPortRepo.On("WithTx", mockTx).Return(mockPortRepo)

	initialUser := &user.User{Id: userID, Balance: startBalance}
	mockUserRepo.On("GetUserForUpdate", mock.Anything, userID).Return(initialUser, nil)

	tradeService := service.NewTradeService(mockUserRepo, mockPortRepo, marketRepo, mockTransactor)
	_, err := tradeService.BuyStock(ctx, userID, symbol, quantity)

	assert.Error(t, err)
	assert.Equal(t, apperrors.ErrInsufficientFunds, err)

	mockTx.AssertNotCalled(t, "Commit", mock.Anything)
	mockTx.AssertCalled(t, "Rollback", mock.Anything)
}
