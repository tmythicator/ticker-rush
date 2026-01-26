package service

import (
	"context"

	"github.com/tmythicator/ticker-rush/server/internal/apperrors"
	"github.com/tmythicator/ticker-rush/server/internal/proto/user"
	valkey "github.com/tmythicator/ticker-rush/server/internal/repository/redis"
)

// TradeService handles stock trading operations.
type TradeService struct {
	userRepo      UserRepository
	portfolioRepo PortfolioRepository
	marketRepo    *valkey.MarketRepository
	transactor    Transactor
}

// NewTradeService creates a new instance of TradeService.
func NewTradeService(
	userRepo UserRepository,
	portfolioRepo PortfolioRepository,
	marketRepo *valkey.MarketRepository,
	transactor Transactor,
) *TradeService {
	return &TradeService{
		userRepo:      userRepo,
		portfolioRepo: portfolioRepo,
		marketRepo:    marketRepo,
		transactor:    transactor,
	}
}

// BuyStock purchases a stock for a user.
func (s *TradeService) BuyStock(
	ctx context.Context,
	userID int64,
	symbol string,
	quantity float64,
) (*user.User, error) {
	// 1. Get current price
	quote, err := s.marketRepo.GetQuote(ctx, symbol)
	if err != nil {
		return nil, err
	}

	cost := quote.GetPrice() * quantity

	// START TRANSACTION
	tx, err := s.transactor.Begin(ctx)
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
	if user.GetBalance() < cost {
		return nil, apperrors.ErrInsufficientFunds
	}

	// 4. Get Current Portfolio Item
	var (
		currentQty float64
		currentAvg float64
	)

	item, err := txPortfolioRepo.GetPortfolioItemForUpdate(ctx, user.GetId(), symbol)
	if err == nil {
		currentQty = item.GetQuantity()
		currentAvg = item.GetAveragePrice()
	}

	// 5. Execute Trade Logic
	user.Balance -= cost

	currentTotalValue := currentQty * currentAvg
	newTotalValue := currentTotalValue + cost
	newTotalQuantity := currentQty + quantity

	newAvgPrice := newTotalValue / newTotalQuantity

	// 6. Persistence
	if err := txUserRepo.SaveUser(ctx, user); err != nil {
		return nil, err
	}

	if err := txPortfolioRepo.SetPortfolioItem(ctx, user.GetId(), symbol, newTotalQuantity, newAvgPrice); err != nil {
		return nil, err
	}

	// COMMIT
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return user, nil
}

// SellStock sells a stock for a user.
func (s *TradeService) SellStock(
	ctx context.Context,
	userID int64,
	symbol string,
	quantity float64,
) (*user.User, error) {
	// 1. Get current price
	quote, err := s.marketRepo.GetQuote(ctx, symbol)
	if err != nil {
		return nil, err
	}

	cost := quote.GetPrice() * quantity

	// START TRANSACTION
	tx, err := s.transactor.Begin(ctx)
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

	if item.GetQuantity() < quantity {
		return nil, apperrors.ErrInsufficientQuantity
	}

	// 3. Get User
	user, err := txUserRepo.GetUserForUpdate(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 4 Execute Trade Logic
	user.Balance += cost
	item.Quantity = item.GetQuantity() - quantity

	// 6. Persistence
	if err := txUserRepo.SaveUser(ctx, user); err != nil {
		return nil, err
	}

	if item.GetQuantity() == 0 {
		err := txPortfolioRepo.DeletePortfolioItem(ctx, userID, symbol)
		if err != nil {
			return nil, err
		}
	} else {
		err := txPortfolioRepo.SetPortfolioItem(ctx, user.GetId(), symbol, item.GetQuantity(), item.GetAveragePrice())
		if err != nil {
			return nil, err
		}
	}

	// COMMIT
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return user, nil
}
