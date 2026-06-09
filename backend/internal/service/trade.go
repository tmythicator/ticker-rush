package service

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
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

	quantityDec := decimal.NewFromFloat(quantity)
	cost := price.Mul(quantityDec)

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

	if balance.LessThan(cost) {
		return nil, apperrors.ErrInsufficientFunds
	}

	// 2. Get Current Portfolio Item
	var (
		currentQty = decimal.Zero
		currentAvg = decimal.Zero
	)

	item, err := txPortfolioRepo.GetPortfolioItemForUpdate(ctx, userID, ladderID, symbol)
	if err == nil {
		currentQty = decimal.NewFromFloat(item.GetQuantity())
		currentAvg = decimal.NewFromFloat(item.GetAveragePrice())
	}

	// 3. Execute Trade Logic
	newBalance := balance.Sub(cost)
	newTotalQuantity := currentQty.Add(quantityDec)
	newAvgPrice := currentQty.Mul(currentAvg).Add(cost).Div(newTotalQuantity)

	// 4. Persistence
	if err := txUserRepo.UpdateUserBalance(ctx, userID, ladderID, newBalance); err != nil {
		return nil, err
	}
	balanceVal, _ := newBalance.Float64()
	user.Balance = balanceVal

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

	quantityDec := decimal.NewFromFloat(quantity)
	totalSaleValue := price.Mul(quantityDec)

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

	itemQtyDec := decimal.NewFromFloat(item.GetQuantity())
	if itemQtyDec.LessThan(quantityDec) {
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
	newBalance := balance.Add(totalSaleValue)
	newQty := itemQtyDec.Sub(quantityDec)

	// 4. Persistence
	if err := txUserRepo.UpdateUserBalance(ctx, userID, ladderID, newBalance); err != nil {
		return nil, err
	}
	balanceVal, _ := newBalance.Float64()
	user.Balance = balanceVal

	itemAvgPriceDec := decimal.NewFromFloat(item.GetAveragePrice())
	if err := s.updatePortfolioPersistence(ctx, txPortfolioRepo, userID, ladderID, symbol, newQty, itemAvgPriceDec); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *TradeService) validateMarketAndParticipation(ctx context.Context, userID int64, symbol string) (decimal.Decimal, int64, error) {
	quote, err := s.marketRepo.GetQuote(ctx, symbol)
	if err != nil {
		return decimal.Zero, 0, err
	}

	if quote.GetIsClosed() {
		return decimal.Zero, 0, apperrors.ErrMarketClosed
	}

	ladderID, err := s.ladderRepo.GetActiveLadder(ctx)
	if err != nil {
		return decimal.Zero, 0, err
	}

	l, err := s.ladderRepo.GetLadder(ctx, ladderID)
	if err != nil {
		return decimal.Zero, 0, err
	}

	now := time.Now().Unix()
	if now < l.GetStartTime().GetSeconds() || now >= l.GetEndTime().GetSeconds() || !l.GetIsActive() {
		return decimal.Zero, 0, apperrors.ErrLadderNotActive
	}

	joined, err := s.ladderRepo.IsUserInLadder(ctx, ladderID, userID)
	if err != nil {
		return decimal.Zero, 0, err
	}
	if !joined {
		return decimal.Zero, 0, apperrors.ErrNotJoinedLadder
	}

	return decimal.NewFromFloat(quote.GetPrice()), ladderID, nil
}

func (s *TradeService) updatePortfolioPersistence(
	ctx context.Context,
	repo PortfolioRepository,
	userID int64,
	ladderID int64,
	symbol string,
	newQty decimal.Decimal,
	avgPrice decimal.Decimal,
) error {
	if newQty.IsZero() {
		return repo.DeletePortfolioItem(ctx, userID, ladderID, symbol)
	}

	return repo.SetPortfolioItem(ctx, userID, ladderID, symbol, newQty, avgPrice)
}
