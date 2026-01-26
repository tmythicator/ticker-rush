// Package mocks provides mock implementations for testing.
package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	pb "github.com/tmythicator/ticker-rush/server/internal/proto/user"
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

// GetUser retrieves a user by ID.
func (m *MockUserRepository) GetUser(ctx context.Context, id int64) (*pb.User, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(*pb.User), args.Error(1)
}

// GetUserByEmail retrieves a user by email.
func (m *MockUserRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*pb.User, string, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}

	return args.Get(0).(*pb.User), args.String(1), args.Error(2)
}

// CreateUser creates a new user.
func (m *MockUserRepository) CreateUser(
	ctx context.Context,
	email string,
	hashedPassword string,
	firstName string,
	lastName string,
	balance float64,
) (*pb.User, error) {
	args := m.Called(ctx, email, hashedPassword, firstName, lastName, balance)

	return args.Get(0).(*pb.User), args.Error(1)
}

// GetUserForUpdate retrieves a user by ID with a lock.
func (m *MockUserRepository) GetUserForUpdate(ctx context.Context, id int64) (*pb.User, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(*pb.User), args.Error(1)
}

// SaveUser updates a user.
func (m *MockUserRepository) SaveUser(ctx context.Context, user *pb.User) error {
	args := m.Called(ctx, user)

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
) ([]*pb.PortfolioItem, error) {
	args := m.Called(ctx, userID)

	return args.Get(0).([]*pb.PortfolioItem), args.Error(1)
}

// GetPortfolioItem retrieves a portfolio item.
func (m *MockPortfolioRepository) GetPortfolioItem(
	ctx context.Context,
	userID int64,
	symbol string,
) (*pb.PortfolioItem, error) {
	args := m.Called(ctx, userID, symbol)

	return args.Get(0).(*pb.PortfolioItem), args.Error(1)
}

// GetPortfolioItemForUpdate retrieves a portfolio item with a lock.
func (m *MockPortfolioRepository) GetPortfolioItemForUpdate(
	ctx context.Context,
	userID int64,
	symbol string,
) (*pb.PortfolioItem, error) {
	args := m.Called(ctx, userID, symbol)

	return args.Get(0).(*pb.PortfolioItem), args.Error(1)
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
