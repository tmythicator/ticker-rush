package grpc

import (
	"context"

	"github.com/tmythicator/ticker-rush/server/internal/api/middleware"
	"github.com/tmythicator/ticker-rush/server/internal/proto/exchange"
	"github.com/tmythicator/ticker-rush/server/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ExchangeServer struct {
	exchange.UnimplementedExchangeServiceServer
	tradeService  *service.TradeService
	marketService *service.MarketService
}

func NewExchangeServer(tradeService *service.TradeService, marketService *service.MarketService) *ExchangeServer {
	return &ExchangeServer{
		tradeService:  tradeService,
		marketService: marketService,
	}
}

func (s *ExchangeServer) GetQuote(ctx context.Context, req *exchange.GetQuoteRequest) (*exchange.GetQuoteResponse, error) {
	_, ok := ctx.Value(middleware.UserIDContextKey).(int64)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user ID not found in context")
	}

	quote, err := s.marketService.GetQuote(ctx, req.Symbol)
	if err != nil {
		return nil, err
	}
	return &exchange.GetQuoteResponse{Quote: quote}, nil
}

func (s *ExchangeServer) BuyStock(ctx context.Context, req *exchange.BuyStockRequest) (*exchange.BuyStockResponse, error) {
	userID, ok := ctx.Value(middleware.UserIDContextKey).(int64)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user ID not found in context")
	}

	_, err := s.tradeService.BuyStock(ctx, userID, req.Symbol, req.Quantity)
	if err != nil {
		return &exchange.BuyStockResponse{Success: false, Message: err.Error()}, nil
	}
	return &exchange.BuyStockResponse{Success: true, Message: "Stock bought successfully"}, nil
}

func (s *ExchangeServer) SellStock(ctx context.Context, req *exchange.SellStockRequest) (*exchange.SellStockResponse, error) {
	userID, ok := ctx.Value(middleware.UserIDContextKey).(int64)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user ID not found in context")
	}

	_, err := s.tradeService.SellStock(ctx, userID, req.Symbol, req.Quantity)
	if err != nil {
		return &exchange.SellStockResponse{Success: false, Message: err.Error()}, nil
	}
	return &exchange.SellStockResponse{Success: true, Message: "Stock sold successfully"}, nil
}
