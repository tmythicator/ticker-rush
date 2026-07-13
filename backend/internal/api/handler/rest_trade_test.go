package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/tmythicator/ticker-rush/backend/internal/apperrors"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/exchange/v1"
	redisRepo "github.com/tmythicator/ticker-rush/backend/internal/repository/redis"
)

func TestBuyStock(t *testing.T) {
	const (
		symbol                  = "AAPL"
		balance         float64 = 10000.0
		price           float64 = 150.0
		quantity        float64 = 2.0
		expectedBalance float64 = balance - price*quantity
	)

	router, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	quote := &redisRepo.ValkeyQuote{Symbol: symbol, Price: price, Timestamp: time.Now().Unix()}
	quoteBytes, _ := json.Marshal(quote)
	valkeyClient.Set(ctx, "market:"+symbol, quoteBytes, 0)

	user, token, activeLadderID := setupJoinedUser(ctx, t, router, balance)

	reqBodyObj := &exchange.CreateTradeRequest{
		Symbol:   symbol,
		Quantity: quantity,
		Action:   exchange.TradeAction_BUY,
	}
	reqBytes, _ := json.Marshal(reqBodyObj)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/trades", bytes.NewReader(reqBytes))
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var rawResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &rawResp)
	assert.NoError(t, err)

	participantMap, ok := rawResp["participant"].(map[string]interface{})
	assert.True(t, ok, "should have participant map")
	userMap, ok := participantMap["user"].(map[string]interface{})
	assert.True(t, ok, "should have user map")

	assert.Equal(t, testUsername, userMap["username"])
	assert.Equal(t, expectedBalance, userMap["balance"])
	assertPublicProfilePrivacy(t, userMap)

	updatedUser, _ := userRepo.GetUser(ctx, user.ID)
	assert.Equal(t, testUsername, updatedUser.Username)
	balanceVal, _ := userRepo.GetUserBalance(ctx, user.ID, activeLadderID)
	assert.Equal(t, expectedBalance, balanceVal.InexactFloat64())

	item, err := portfolioRepo.GetPortfolioItem(ctx, user.ID, activeLadderID, symbol)
	assert.NoError(t, err)
	assert.Equal(t, quantity, item.Quantity.InexactFloat64())
}

func TestSellStock(t *testing.T) {
	router, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	const (
		mockPrice                 float64 = 150.0
		mockStartBalance          float64 = 10000.0
		mockPortfolioQuantity     float64 = 5.0
		mockSellQuantity          float64 = 2.0
		expectedPortfolioQuantity float64 = mockPortfolioQuantity - mockSellQuantity
		expectedBalance           float64 = mockPrice*mockSellQuantity + mockStartBalance
	)

	quote := &redisRepo.ValkeyQuote{Symbol: "AAPL", Price: mockPrice, Timestamp: time.Now().Unix()}
	quoteBytes, _ := json.Marshal(quote)
	valkeyClient.Set(ctx, "market:AAPL", quoteBytes, 0)

	user, token, activeLadderID := setupJoinedUser(ctx, t, router, mockStartBalance)

	err := portfolioRepo.SetPortfolioItem(
		ctx,
		user.ID,
		activeLadderID,
		"AAPL",
		decimal.NewFromFloat(mockPortfolioQuantity),
		decimal.NewFromFloat(mockPrice),
	)
	assert.NoError(t, err)

	reqBodyObj := &exchange.CreateTradeRequest{
		Symbol:   "AAPL",
		Quantity: mockSellQuantity,
		Action:   exchange.TradeAction_SELL,
	}
	reqBytes, _ := json.Marshal(reqBodyObj)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/trades", bytes.NewReader(reqBytes))
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var rawResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &rawResp)
	assert.NoError(t, err)

	participantMap, ok := rawResp["participant"].(map[string]interface{})
	assert.True(t, ok, "should have participant map")
	userMap, ok := participantMap["user"].(map[string]interface{})
	assert.True(t, ok, "should have user map")

	assert.Equal(t, testUsername, userMap["username"])
	assert.Equal(t, expectedBalance, userMap["balance"])
	assertPublicProfilePrivacy(t, userMap)

	balanceVal, _ := userRepo.GetUserBalance(ctx, user.ID, activeLadderID)
	assert.Equal(t, expectedBalance, balanceVal.InexactFloat64())

	item, err := portfolioRepo.GetPortfolioItem(ctx, user.ID, activeLadderID, "AAPL")
	assert.NoError(t, err)
	assert.Equal(t, expectedPortfolioQuantity, item.Quantity.InexactFloat64())
}

func TestInsufficientFunds(t *testing.T) {
	router, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	const (
		mockPrice        = 151.0
		mockStartBalance = 20.0
		mockBuyQuantity  = 1
	)

	quote := &redisRepo.ValkeyQuote{Symbol: "AAPL", Price: mockPrice, Timestamp: time.Now().Unix()}
	quoteBytes, _ := json.Marshal(quote)
	valkeyClient.Set(ctx, "market:AAPL", quoteBytes, 0)

	_, token, _ := setupJoinedUser(ctx, t, router, mockStartBalance)

	reqBodyObj := &exchange.CreateTradeRequest{
		Symbol:   "AAPL",
		Quantity: float64(mockBuyQuantity),
		Action:   exchange.TradeAction_BUY,
	}
	reqBytes, _ := json.Marshal(reqBodyObj)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/trades", bytes.NewReader(reqBytes))
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusPaymentRequired, w.Code)
	assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))
	var prob apperrors.ProblemDetails
	err := json.Unmarshal(w.Body.Bytes(), &prob)
	assert.NoError(t, err)
	assert.Equal(t, apperrors.TypeInsufficientFunds, prob.Type)
	assert.Equal(t, apperrors.ErrInsufficientFunds.Error(), prob.Detail)
}

func TestSellAllStock(t *testing.T) {
	const (
		symbol                      = "AAPL"
		mockStartBalance    float64 = 0.0
		mockPrice           float64 = 150.0
		mockQuantity        float64 = 5.0
		mockSellQuantity    float64 = 5.0
		mockExpectedBalance float64 = mockStartBalance + mockSellQuantity*mockPrice
	)

	router, mr, pool := setupTestRouter(t)
	defer mr.Close()
	defer pool.Close()

	quote := &redisRepo.ValkeyQuote{Symbol: symbol, Price: mockPrice, Timestamp: time.Now().Unix()}
	quoteBytes, _ := json.Marshal(quote)
	valkeyClient.Set(ctx, "market:"+symbol, quoteBytes, 0)

	user, token, activeLadderID := setupJoinedUser(ctx, t, router, mockStartBalance)

	err := portfolioRepo.SetPortfolioItem(ctx, user.ID, activeLadderID, symbol, decimal.NewFromFloat(mockQuantity), decimal.NewFromFloat(mockPrice))
	assert.NoError(t, err)

	reqBodyObj := &exchange.CreateTradeRequest{
		Symbol:   symbol,
		Quantity: mockSellQuantity,
		Action:   exchange.TradeAction_SELL,
	}
	reqBytes, _ := json.Marshal(reqBodyObj)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/trades", bytes.NewReader(reqBytes))
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	balanceVal, _ := userRepo.GetUserBalance(ctx, user.ID, activeLadderID)
	assert.Equal(t, mockExpectedBalance, balanceVal.InexactFloat64())

	_, err = portfolioRepo.GetPortfolioItem(ctx, user.ID, activeLadderID, symbol)
	assert.Error(t, err)
}
