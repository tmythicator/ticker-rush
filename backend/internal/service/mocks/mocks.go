// Package mocks provides mock implementations for testing.
package mocks

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/tmythicator/ticker-rush/server/internal/proto/exchange/v1"
	"github.com/tmythicator/ticker-rush/server/internal/proto/user/v1"
	"github.com/tmythicator/ticker-rush/server/internal/service"
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
func (m *MockUserRepository) GetUsers(ctx context.Context) ([]*user.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*user.User), args.Error(1)
}

// GetUser retrieves a user by ID.
func (m *MockUserRepository) GetUser(ctx context.Context, id int64) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*user.User), args.Error(1)
}

// GetUserByUsername retrieves a user by username.
func (m *MockUserRepository) GetUserByUsername(
	ctx context.Context,
	username string,
) (*user.User, string, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}

	return args.Get(0).(*user.User), args.String(1), args.Error(2)
}

// CreateUser creates a new user.
func (m *MockUserRepository) CreateUser(
	ctx context.Context,
	username string,
	hashedPassword string,
	firstName string,
	lastName string,
	balance float64,
	website string,
	isPublic bool,
) (*user.User, error) {
	args := m.Called(ctx, username, hashedPassword, firstName, lastName, balance, website, isPublic)

	return args.Get(0).(*user.User), args.Error(1)
}

// GetUserForUpdate retrieves a user by ID with a lock.
func (m *MockUserRepository) GetUserForUpdate(ctx context.Context, id int64) (*user.User, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(*user.User), args.Error(1)
}

// UpdateUserProfile updates a user's profile.
func (m *MockUserRepository) UpdateUserProfile(ctx context.Context, user *user.User) error {
	args := m.Called(ctx, user)

	return args.Error(0)
}

// UpdateUserBalance updates the user's balance.
func (m *MockUserRepository) UpdateUserBalance(ctx context.Context, id int64, balance float64) error {
	args := m.Called(ctx, id, balance)

	return args.Error(0)
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
) ([]*user.PortfolioItem, error) {
	args := m.Called(ctx, userID)

	return args.Get(0).([]*user.PortfolioItem), args.Error(1)
}

// GetPortfolioItem retrieves a portfolio item.
func (m *MockPortfolioRepository) GetPortfolioItem(
	ctx context.Context,
	userID int64,
	symbol string,
) (*user.PortfolioItem, error) {
	args := m.Called(ctx, userID, symbol)

	return args.Get(0).(*user.PortfolioItem), args.Error(1)
}

// GetPortfolioItemForUpdate retrieves a portfolio item with a lock.
func (m *MockPortfolioRepository) GetPortfolioItemForUpdate(
	ctx context.Context,
	userID int64,
	symbol string,
) (*user.PortfolioItem, error) {
	args := m.Called(ctx, userID, symbol)

	return args.Get(0).(*user.PortfolioItem), args.Error(1)
}

// SetPortfolioItem updates or inserts a portfolio item.
func (m *MockPortfolioRepository) SetPortfolioItem(
	ctx context.Context,
	userID int64,
	symbol string,
	quantity float64,
	averagePrice float64,
) error {
	args := m.Called(ctx, userID, symbol, quantity, averagePrice)

	return args.Error(0)
}

// DeletePortfolioItem removes a portfolio item.
func (m *MockPortfolioRepository) DeletePortfolioItem(
	ctx context.Context,
	userID int64,
	symbol string,
) error {
	args := m.Called(ctx, userID, symbol)

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
func (m *MockMarketRepository) GetQuote(ctx context.Context, symbol string) (*exchange.Quote, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*exchange.Quote), args.Error(1)
}

// SaveQuote saves a stock quote.
func (m *MockMarketRepository) SaveQuote(ctx context.Context, quote *exchange.Quote) error {
	args := m.Called(ctx, quote)

	return args.Error(0)
}

// SubscribeToQuotes subscribes to quotes.
func (m *MockMarketRepository) SubscribeToQuotes(ctx context.Context, symbol string) *redis.PubSub {
	args := m.Called(ctx, symbol)

	return args.Get(0).(*redis.PubSub)
}
