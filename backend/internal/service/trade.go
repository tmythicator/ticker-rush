package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tmythicator/ticker-rush/server/internal/repository/postgres"
	valkey "github.com/tmythicator/ticker-rush/server/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/server/model"
	pb "github.com/tmythicator/ticker-rush/server/proto/user"
)

type TradeService struct {
	userRepo      *postgres.UserRepository
	portfolioRepo *postgres.PortfolioRepository
	marketRepo    *valkey.MarketRepository
	postgresPool  *pgxpool.Pool
}

func NewTradeService(userRepo *postgres.UserRepository, portfolioRepo *postgres.PortfolioRepository, marketRepo *valkey.MarketRepository, postgresPool *pgxpool.Pool) *TradeService {
	return &TradeService{
		userRepo:      userRepo,
		portfolioRepo: portfolioRepo,
		marketRepo:    marketRepo,
		postgresPool:  postgresPool,
	}
}

func (s *TradeService) BuyStock(ctx context.Context, userID int64, symbol string, quantity int32) (*pb.User, error) {
	// 1. Get current price
	quote, err := s.marketRepo.GetQuote(ctx, symbol)
	if err != nil {
		return nil, err
	}

	cost := quote.Price * float64(quantity)

	// START TRANSACTION
	tx, err := s.postgresPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Create Transactional Repos
	txUserRepo := s.userRepo.WithTx(tx)
	txPortfolioRepo := s.portfolioRepo.WithTx(tx)

	// 2. Get User
	user, err := txUserRepo.GetUserForUpdate(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 3. Check Balance
	if user.Balance < cost {
		return nil, model.ErrInsufficientFunds
	}

	// 4. Get Current Portfolio Item
	var currentQty int32 = 0
	var currentAvg float64 = 0
	item, err := txPortfolioRepo.GetPortfolioItemForUpdate(ctx, user.Id, symbol)
	if err == nil {
		currentQty = item.Quantity
		currentAvg = item.AveragePrice
	}

	// 5. Execute Trade Logic
	user.Balance -= cost

	currentTotalValue := float64(currentQty) * currentAvg
	newTotalValue := currentTotalValue + cost
	newTotalQuantity := currentQty + quantity

	item.AveragePrice = newTotalValue / float64(newTotalQuantity)
	item.Quantity = newTotalQuantity

	// 6. Persistence
	if err := txUserRepo.SaveUser(ctx, user); err != nil {
		return nil, err
	}

	if err := txPortfolioRepo.SetPortfolioItem(ctx, user.Id, symbol, item.Quantity, item.AveragePrice); err != nil {
		return nil, err
	}

	// COMMIT
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *TradeService) SellStock(ctx context.Context, userID int64, symbol string, quantity int32) (*pb.User, error) {
	// 1. Get current price
	quote, err := s.marketRepo.GetQuote(ctx, symbol)
	if err != nil {
		return nil, err
	}

	cost := quote.Price * float64(quantity)

	// START TRANSACTION
	tx, err := s.postgresPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Create Transactional Repos
	txUserRepo := s.userRepo.WithTx(tx)
	txPortfolioRepo := s.portfolioRepo.WithTx(tx)

	// 2. Check Portfolio Item
	item, err := txPortfolioRepo.GetPortfolioItemForUpdate(ctx, userID, symbol)
	if err != nil {
		return nil, err
	}

	if item.Quantity < quantity {
		return nil, model.ErrInsufficientQuantity
	}

	// 3. Get User
	user, err := txUserRepo.GetUserForUpdate(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 4 Execute Trade Logic
	user.Balance += cost
	item.Quantity = item.Quantity - quantity

	// 6. Persistence
	if err := txUserRepo.SaveUser(ctx, user); err != nil {
		return nil, err
	}

	if item.Quantity == 0 {
		if err := txPortfolioRepo.DeletePortfolioItem(ctx, userID, symbol); err != nil {
			return nil, err
		}
	} else {
		if err := txPortfolioRepo.SetPortfolioItem(ctx, user.Id, symbol, item.Quantity, item.AveragePrice); err != nil {
			return nil, err
		}
	}

	// COMMIT
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return user, nil
}
