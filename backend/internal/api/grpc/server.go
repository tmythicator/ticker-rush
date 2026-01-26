// Package grpc provides the gRPC server implementation.
package grpc

import (
	"context"

	"github.com/tmythicator/ticker-rush/server/internal/api/middleware"
	"github.com/tmythicator/ticker-rush/server/internal/proto/exchange"
	"github.com/tmythicator/ticker-rush/server/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ExchangeServer implements the gRPC exchange service.
type ExchangeServer struct {
	exchange.UnimplementedExchangeServiceServer

	tradeService  *service.TradeService
	marketService *service.MarketService
}

// NewExchangeServer creates a new instance of ExchangeServer.
func NewExchangeServer(
	tradeService *service.TradeService,
	marketService *service.MarketService,
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
	_, ok := ctx.Value(middleware.UserIDContextKey).(int64)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user ID not found in context")
	}

	quote, err := s.marketService.GetQuote(ctx, req.GetSymbol())
	if err != nil {
		return nil, err
	}

	return &exchange.GetQuoteResponse{Quote: quote}, nil
}

// BuyStock executes a buy order.
func (s *ExchangeServer) BuyStock(
	ctx context.Context,
	req *exchange.BuyStockRequest,
) (*exchange.BuyStockResponse, error) {
	userID, ok := ctx.Value(middleware.UserIDContextKey).(int64)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user ID not found in context")
	}

	_, err := s.tradeService.BuyStock(ctx, userID, req.GetSymbol(), req.GetQuantity())
	if err != nil {
		return &exchange.BuyStockResponse{Success: false, Message: err.Error()}, nil
	}

	return &exchange.BuyStockResponse{Success: true, Message: "Stock bought successfully"}, nil
}

// SellStock executes a sell order.
func (s *ExchangeServer) SellStock(
	ctx context.Context,
	req *exchange.SellStockRequest,
) (*exchange.SellStockResponse, error) {
	userID, ok := ctx.Value(middleware.UserIDContextKey).(int64)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user ID not found in context")
	}

	_, err := s.tradeService.SellStock(ctx, userID, req.GetSymbol(), req.GetQuantity())
	if err != nil {
		return &exchange.SellStockResponse{Success: false, Message: err.Error()}, nil
	}

	return &exchange.SellStockResponse{Success: true, Message: "Stock sold successfully"}, nil
}
