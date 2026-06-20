// Package mocks provides mock implementations for testing.
package mocks

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"

	"github.com/tmythicator/ticker-rush/backend/internal/domain"
	"github.com/tmythicator/ticker-rush/backend/internal/service"
)

// MockTransaction is a mock implementation of Transaction.
type MockTransaction struct {
	mock.Mock
}

// Commit commits the transaction.
func (m *MockTransaction) Commit(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}

// Rollback rolls back the transaction.
func (m *MockTransaction) Rollback(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}

// MockTransactor is a mock implementation of Transactor.
type MockTransactor struct {
	mock.Mock
}

// Begin starts a new mock transaction.
func (m *MockTransactor) Begin(ctx context.Context) (service.Transaction, error) {
	args := m.Called(ctx)

	return args.Get(0).(service.Transaction), args.Error(1)
}

// MockUserRepository is a mock implementation of UserRepository.
type MockUserRepository struct {
	mock.Mock
}

// GetUsers retrieves all users
func (m *MockUserRepository) GetUsers(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.User), args.Error(1)
}

// GetUser retrieves a user by ID.
func (m *MockUserRepository) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.User), args.Error(1)
}

// GetUserByUsername retrieves a user by username.
func (m *MockUserRepository) GetUserByUsername(
	ctx context.Context,
	username string,
) (*domain.User, string, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}

	return args.Get(0).(*domain.User), args.String(1), args.Error(2)
}

// CreateUser creates a new user.
func (m *MockUserRepository) CreateUser(
	ctx context.Context,
	params service.CreateUserParams,
) (*domain.User, error) {
	args := m.Called(ctx, params)

	return args.Get(0).(*domain.User), args.Error(1)
}

// GetUserForUpdate retrieves a user by ID with a lock.
func (m *MockUserRepository) GetUserForUpdate(ctx context.Context, id int64) (*domain.User, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(*domain.User), args.Error(1)
}

// UpdateUserProfile updates a user's profile.
func (m *MockUserRepository) UpdateUserProfile(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)

	return args.Error(0)
}

// UpdateUserBalance updates the user's balance.
func (m *MockUserRepository) UpdateUserBalance(ctx context.Context, id int64, ladderID int64, balance decimal.Decimal) error {
	args := m.Called(ctx, id, ladderID, balance)

	return args.Error(0)
}

// GetUserBalance retrieves the user's balance.
func (m *MockUserRepository) GetUserBalance(ctx context.Context, userID int64, ladderID int64) (decimal.Decimal, error) {
	args := m.Called(ctx, userID, ladderID)

	return args.Get(0).(decimal.Decimal), args.Error(1)
}

// GetUserWithPortfolioForActiveLadder retrieves a user with their portfolio items.
func (m *MockUserRepository) GetUserWithPortfolioForActiveLadder(ctx context.Context, userID int64) (*domain.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.User), args.Error(1)
}

// WithTx returns a new UserRepository with the transaction.
func (m *MockUserRepository) WithTx(tx service.Transaction) service.UserRepository {
	args := m.Called(tx)

	return args.Get(0).(service.UserRepository)
}

// MockPortfolioRepository is a mock implementation of PortfolioRepository.
type MockPortfolioRepository struct {
	mock.Mock
}

// GetPortfolio retrieves the portfolio.
func (m *MockPortfolioRepository) GetPortfolio(
	ctx context.Context,
	userID int64,
	ladderID int64,
) ([]*domain.PortfolioItem, error) {
	args := m.Called(ctx, userID, ladderID)

	return args.Get(0).([]*domain.PortfolioItem), args.Error(1)
}

// GetPortfolioItem retrieves a portfolio item.
func (m *MockPortfolioRepository) GetPortfolioItem(
	ctx context.Context,
	userID int64,
	ladderID int64,
	symbol string,
) (*domain.PortfolioItem, error) {
	args := m.Called(ctx, userID, ladderID, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.PortfolioItem), args.Error(1)
}

// GetPortfolioItemForUpdate retrieves a portfolio item with a lock.
func (m *MockPortfolioRepository) GetPortfolioItemForUpdate(
	ctx context.Context,
	userID int64,
	ladderID int64,
	symbol string,
) (*domain.PortfolioItem, error) {
	args := m.Called(ctx, userID, ladderID, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.PortfolioItem), args.Error(1)
}

// SetPortfolioItem updates or inserts a portfolio item.
func (m *MockPortfolioRepository) SetPortfolioItem(
	ctx context.Context,
	userID int64,
	ladderID int64,
	symbol string,
	quantity decimal.Decimal,
	averagePrice decimal.Decimal,
) error {
	args := m.Called(ctx, userID, ladderID, symbol, quantity, averagePrice)

	return args.Error(0)
}

// DeletePortfolioItem removes a portfolio item.
func (m *MockPortfolioRepository) DeletePortfolioItem(
	ctx context.Context,
	userID int64,
	ladderID int64,
	symbol string,
) error {
	args := m.Called(ctx, userID, ladderID, symbol)

	return args.Error(0)
}

// WithTx returns a new PortfolioRepository with the transaction.
func (m *MockPortfolioRepository) WithTx(tx service.Transaction) service.PortfolioRepository {
	args := m.Called(tx)

	return args.Get(0).(service.PortfolioRepository)
}

// MockMarketRepository is a mock implementation of MarketRepository.
type MockMarketRepository struct {
	mock.Mock
}

// GetQuote retrieves a stock quote.
func (m *MockMarketRepository) GetQuote(ctx context.Context, symbol string) (*domain.Quote, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Quote), args.Error(1)
}

// SaveQuote saves a stock quote.
func (m *MockMarketRepository) SaveQuote(ctx context.Context, quote *domain.Quote) error {
	args := m.Called(ctx, quote)

	return args.Error(0)
}

// SubscribeToQuotes subscribes to quotes.
func (m *MockMarketRepository) SubscribeToQuotes(ctx context.Context, symbol string) *redis.PubSub {
	args := m.Called(ctx, symbol)

	return args.Get(0).(*redis.PubSub)
}

// MockLadderRepository is a mock implementation of LadderRepository.
type MockLadderRepository struct {
	mock.Mock
}

// GetActiveLadder retrieves the currently active ladder ID.
func (m *MockLadderRepository) GetActiveLadder(ctx context.Context) (int64, error) {
	args := m.Called(ctx)

	return args.Get(0).(int64), args.Error(1)
}

// GetLadder retrieves a ladder by ID.
func (m *MockLadderRepository) GetLadder(ctx context.Context, id int64) (*domain.Ladder, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Ladder), args.Error(1)
}

// GetAllowedTickers retrieves the allowed stock symbols for a given ladder.
func (m *MockLadderRepository) GetAllowedTickers(ctx context.Context, ladderID int64) ([]*domain.TickerInfo, error) {
	args := m.Called(ctx, ladderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.TickerInfo), args.Error(1)
}

// JoinLadder mock implementation.
func (m *MockLadderRepository) JoinLadder(ctx context.Context, ladderID int64, userID int64) error {
	args := m.Called(ctx, ladderID, userID)

	return args.Error(0)
}

// IsUserInLadder mock.
func (m *MockLadderRepository) IsUserInLadder(ctx context.Context, ladderID int64, userID int64) (bool, error) {
	args := m.Called(ctx, ladderID, userID)

	return args.Bool(0), args.Error(1)
}

// GetExpiredActiveLadders mock.
func (m *MockLadderRepository) GetExpiredActiveLadders(ctx context.Context, now time.Time) ([]*domain.Ladder, error) {
	args := m.Called(ctx, now)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.Ladder), args.Error(1)
}

// GetPendingLaddersToActivate mock.
func (m *MockLadderRepository) GetPendingLaddersToActivate(ctx context.Context, now time.Time) ([]*domain.Ladder, error) {
	args := m.Called(ctx, now)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.Ladder), args.Error(1)
}

// UpdateLadderStatus mock.
func (m *MockLadderRepository) UpdateLadderStatus(ctx context.Context, id int64, isActive bool) error {
	args := m.Called(ctx, id, isActive)

	return args.Error(0)
}

// GetLadderParticipants mock.
func (m *MockLadderRepository) GetLadderParticipants(ctx context.Context, ladderID int64) ([]domain.LadderParticipant, error) {
	args := m.Called(ctx, ladderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]domain.LadderParticipant), args.Error(1)
}

// InsertLadderParticipant mock.
func (m *MockLadderRepository) InsertLadderParticipant(ctx context.Context, ladderID int64, userID int64, finalBalance decimal.Decimal, finalRank int32) error {
	args := m.Called(ctx, ladderID, userID, finalBalance, finalRank)

	return args.Error(0)
}

// PruneLadderParticipants mock.
func (m *MockLadderRepository) PruneLadderParticipants(ctx context.Context, ladderID int64, rankThreshold int32) error {
	args := m.Called(ctx, ladderID, rankThreshold)

	return args.Error(0)
}

// DeleteLadderPortfolioItemsByLadder mock.
func (m *MockLadderRepository) DeleteLadderPortfolioItemsByLadder(ctx context.Context, ladderID int64) error {
	args := m.Called(ctx, ladderID)

	return args.Error(0)
}

// MockLeaderboardRepository is a mock implementation of LeaderboardRepository.
type MockLeaderboardRepository struct {
	mock.Mock
}

// UpdateRank mock.
func (m *MockLeaderboardRepository) UpdateRank(ctx context.Context, ladderID int64, userID int64, score float64) error {
	args := m.Called(ctx, ladderID, userID, score)

	return args.Error(0)
}

// GetLeaderboard mock.
func (m *MockLeaderboardRepository) GetLeaderboard(ctx context.Context, ladderID int64, offset int, limit int) ([]service.LeaderboardScore, error) {
	args := m.Called(ctx, ladderID, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]service.LeaderboardScore), args.Error(1)
}

// GetTotalCount mock.
func (m *MockLeaderboardRepository) GetTotalCount(ctx context.Context, ladderID int64) (int64, error) {
	args := m.Called(ctx, ladderID)

	return int64(args.Int(0)), args.Error(1)
}

// GetLastUpdate mock.
func (m *MockLeaderboardRepository) GetLastUpdate(ctx context.Context, ladderID int64) (int64, error) {
	args := m.Called(ctx, ladderID)

	return args.Get(0).(int64), args.Error(1)
}

// SetLastUpdate mock.
func (m *MockLeaderboardRepository) SetLastUpdate(ctx context.Context, ladderID int64, timestamp int64) error {
	args := m.Called(ctx, ladderID, timestamp)

	return args.Error(0)
}

// RemoveFromLeaderboard mock.
func (m *MockLeaderboardRepository) RemoveFromLeaderboard(ctx context.Context, ladderID int64, userID int64) error {
	args := m.Called(ctx, ladderID, userID)

	return args.Error(0)
}

// MockHistoryRepository is a mock implementation of HistoryRepository.
type MockHistoryRepository struct {
	mock.Mock
}

// SaveQuote mock.
func (m *MockHistoryRepository) SaveQuote(ctx context.Context, quote *domain.Quote) error {
	args := m.Called(ctx, quote)

	return args.Error(0)
}

// GetHistory mock.
func (m *MockHistoryRepository) GetHistory(ctx context.Context, symbol string, limit int) ([]*domain.Quote, error) {
	args := m.Called(ctx, symbol, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.Quote), args.Error(1)
}
