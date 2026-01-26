package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	pb "github.com/tmythicator/ticker-rush/server/internal/proto/user"
	"github.com/tmythicator/ticker-rush/server/internal/service"
)

type MockTransaction struct {
	mock.Mock
}

func (m *MockTransaction) Commit(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}

func (m *MockTransaction) Rollback(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}

type MockTransactor struct {
	mock.Mock
}

func (m *MockTransactor) Begin(ctx context.Context) (service.Transaction, error) {
	args := m.Called(ctx)

	return args.Get(0).(service.Transaction), args.Error(1)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUser(ctx context.Context, id int64) (*pb.User, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(*pb.User), args.Error(1)
}

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

func (m *MockUserRepository) GetUserForUpdate(ctx context.Context, id int64) (*pb.User, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(*pb.User), args.Error(1)
}

func (m *MockUserRepository) SaveUser(ctx context.Context, user *pb.User) error {
	args := m.Called(ctx, user)

	return args.Error(0)
}

func (m *MockUserRepository) WithTx(tx service.Transaction) service.UserRepository {
	args := m.Called(tx)

	return args.Get(0).(service.UserRepository)
}

type MockPortfolioRepository struct {
	mock.Mock
}

func (m *MockPortfolioRepository) GetPortfolio(
	ctx context.Context,
	userID int64,
) ([]*pb.PortfolioItem, error) {
	args := m.Called(ctx, userID)

	return args.Get(0).([]*pb.PortfolioItem), args.Error(1)
}

func (m *MockPortfolioRepository) GetPortfolioItem(
	ctx context.Context,
	userID int64,
	symbol string,
) (*pb.PortfolioItem, error) {
	args := m.Called(ctx, userID, symbol)

	return args.Get(0).(*pb.PortfolioItem), args.Error(1)
}

func (m *MockPortfolioRepository) GetPortfolioItemForUpdate(
	ctx context.Context,
	userID int64,
	symbol string,
) (*pb.PortfolioItem, error) {
	args := m.Called(ctx, userID, symbol)

	return args.Get(0).(*pb.PortfolioItem), args.Error(1)
}

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

func (m *MockPortfolioRepository) DeletePortfolioItem(
	ctx context.Context,
	userID int64,
	symbol string,
) error {
	args := m.Called(ctx, userID, symbol)

	return args.Error(0)
}

func (m *MockPortfolioRepository) WithTx(tx service.Transaction) service.PortfolioRepository {
	args := m.Called(tx)

	return args.Get(0).(service.PortfolioRepository)
}
