package service

import (
	"context"

	"github.com/tmythicator/ticker-rush/backend/internal/apperrors"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/user/v1"
)

// TradeService handles stock trading operations.
type TradeService struct {
	userRepo      UserRepository
	portfolioRepo PortfolioRepository
	marketRepo    MarketRepository
	ladderRepo    LadderRepository
	transactor    Transactor
}

// NewTradeService creates a new instance of TradeService.
func NewTradeService(
	userRepo UserRepository,
	portfolioRepo PortfolioRepository,
	marketRepo MarketRepository,
	ladderRepo LadderRepository,
	transactor Transactor,
) *TradeService {
	return &TradeService{
		userRepo:      userRepo,
		portfolioRepo: portfolioRepo,
		marketRepo:    marketRepo,
		ladderRepo:    ladderRepo,
		transactor:    transactor,
	}
}

// BuyStock purchases a stock for a user for the active ladder.
func (s *TradeService) BuyStock(
	ctx context.Context,
	userID int64,
	symbol string,
	quantity float64,
) (*user.User, error) {
	price, ladderID, err := s.validateMarketAndParticipation(ctx, userID, symbol)
	if err != nil {
		return nil, err
	}

	cost := price * quantity

	// START TRANSACTION
	tx, err := s.transactor.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	txUserRepo := s.userRepo.WithTx(tx)
	txPortfolioRepo := s.portfolioRepo.WithTx(tx)

	// 1. Get User & Balance
	user, err := txUserRepo.GetUserForUpdate(ctx, userID)
	if err != nil {
		return nil, err
	}

	balance, err := txUserRepo.GetUserBalance(ctx, userID, ladderID)
	if err != nil {
		return nil, err
	}

	if balance < cost {
		return nil, apperrors.ErrInsufficientFunds
	}

	// 2. Get Current Portfolio Item
	var (
		currentQty float64
		currentAvg float64
	)

	item, err := txPortfolioRepo.GetPortfolioItemForUpdate(ctx, userID, ladderID, symbol)
	if err == nil {
		currentQty = item.GetQuantity()
		currentAvg = item.GetAveragePrice()
	}

	// 3. Execute Trade Logic
	newBalance := balance - cost
	newTotalQuantity := currentQty + quantity
	newAvgPrice := ((currentQty * currentAvg) + cost) / newTotalQuantity

	// 4. Persistence
	if err := txUserRepo.UpdateUserBalance(ctx, userID, ladderID, newBalance); err != nil {
		return nil, err
	}
	user.Balance = newBalance

	if err := s.updatePortfolioPersistence(ctx, txPortfolioRepo, userID, ladderID, symbol, newTotalQuantity, newAvgPrice); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return user, nil
}

// SellStock sells a stock for a user for the active ladder.
func (s *TradeService) SellStock(
	ctx context.Context,
	userID int64,
	symbol string,
	quantity float64,
) (*user.User, error) {
	price, ladderID, err := s.validateMarketAndParticipation(ctx, userID, symbol)
	if err != nil {
		return nil, err
	}

	totalSaleValue := price * quantity

	// START TRANSACTION
	tx, err := s.transactor.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	txUserRepo := s.userRepo.WithTx(tx)
	txPortfolioRepo := s.portfolioRepo.WithTx(tx)

	// 1. Check Portfolio Item
	item, err := txPortfolioRepo.GetPortfolioItemForUpdate(ctx, userID, ladderID, symbol)
	if err != nil {
		return nil, err
	}

	if item.GetQuantity() < quantity {
		return nil, apperrors.ErrInsufficientQuantity
	}

	// 2. Get User & Balance
	user, err := txUserRepo.GetUserForUpdate(ctx, userID)
	if err != nil {
		return nil, err
	}

	balance, err := txUserRepo.GetUserBalance(ctx, userID, ladderID)
	if err != nil {
		return nil, err
	}

	// 3. Execute Trade Logic
	newBalance := balance + totalSaleValue
	newQty := item.GetQuantity() - quantity

	// 4. Persistence
	if err := txUserRepo.UpdateUserBalance(ctx, userID, ladderID, newBalance); err != nil {
		return nil, err
	}
	user.Balance = newBalance

	if err := s.updatePortfolioPersistence(ctx, txPortfolioRepo, userID, ladderID, symbol, newQty, item.GetAveragePrice()); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *TradeService) validateMarketAndParticipation(ctx context.Context, userID int64, symbol string) (float64, int64, error) {
	quote, err := s.marketRepo.GetQuote(ctx, symbol)
	if err != nil {
		return 0, 0, err
	}

	if quote.GetIsClosed() {
		return 0, 0, apperrors.ErrMarketClosed
	}

	ladderID, err := s.ladderRepo.GetActiveLadder(ctx)
	if err != nil {
		return 0, 0, err
	}

	joined, err := s.ladderRepo.IsUserInLadder(ctx, ladderID, userID)
	if err != nil {
		return 0, 0, err
	}
	if !joined {
		return 0, 0, apperrors.ErrNotJoinedLadder
	}

	return quote.GetPrice(), ladderID, nil
}

func (s *TradeService) updatePortfolioPersistence(
	ctx context.Context,
	repo PortfolioRepository,
	userID int64,
	ladderID int64,
	symbol string,
	newQty float64,
	avgPrice float64,
) error {
	if newQty == 0 {
		return repo.DeletePortfolioItem(ctx, userID, ladderID, symbol)
	}

	return repo.SetPortfolioItem(ctx, userID, ladderID, symbol, newQty, avgPrice)
}
