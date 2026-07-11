// Package grpc provides the gRPC server implementation.
package grpc

import (
	"context"

	"github.com/tmythicator/ticker-rush/backend/internal/api/handler"
	"github.com/tmythicator/ticker-rush/backend/internal/api/middleware"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/exchange/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/service"
)

// ExchangeServer implements the gRPC exchange service.
type ExchangeServer struct {
	exchange.UnimplementedExchangeServiceServer

	tradeService  *service.Trade
	marketService *service.Market
}

// NewExchangeServer creates a new instance of ExchangeServer.
func NewExchangeServer(
	tradeService *service.Trade,
	marketService *service.Market,
) *ExchangeServer {
	return &ExchangeServer{
		tradeService:  tradeService,
		marketService: marketService,
	}
}

// GetQuote retrieves the current price of a stock.
func (s *ExchangeServer) GetQuote(
	ctx context.Context,
	req *exchange.GetQuoteRequest,
) (*exchange.GetQuoteResponse, error) {
	_, err := middleware.GetRequiredUserID(ctx)
	if err != nil {
		return nil, err
	}

	quote, err := s.marketService.GetQuote(ctx, req.GetSymbol())
	if err != nil {
		return nil, err
	}

	return &exchange.GetQuoteResponse{Quote: handler.ToExternalQuote(quote)}, nil
}

// CreateTrade executes a buy or sell order.
func (s *ExchangeServer) CreateTrade(
	ctx context.Context,
	req *exchange.CreateTradeRequest,
) (*exchange.CreateTradeResponse, error) {
	userID, err := middleware.GetRequiredUserID(ctx)
	if err != nil {
		return nil, err
	}

	var msg string
	if req.GetAction() == exchange.TradeAction_BUY {
		_, err = s.tradeService.BuyStock(ctx, userID, req.GetSymbol(), req.GetQuantity())
		msg = "Stock bought successfully"
	} else if req.GetAction() == exchange.TradeAction_SELL {
		_, err = s.tradeService.SellStock(ctx, userID, req.GetSymbol(), req.GetQuantity())
		msg = "Stock sold successfully"
	} else {
		return &exchange.CreateTradeResponse{Success: false, Message: "Invalid trade action"}, nil
	}

	if err != nil {
		return &exchange.CreateTradeResponse{Success: false, Message: err.Error()}, nil
	}

	return &exchange.CreateTradeResponse{Success: true, Message: msg}, nil
}
