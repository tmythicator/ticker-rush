package service

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/tmythicator/ticker-rush/server/internal/proto/exchange/v1"
	"github.com/tmythicator/ticker-rush/server/internal/proto/user/v1"
)

// Transaction represents a database transaction.
type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// Transactor is responsible for creating transactions.
type Transactor interface {
	Begin(ctx context.Context) (Transaction, error)
}

// MarketRepository defines the interface for market data persistence.
type MarketRepository interface {
	GetQuote(ctx context.Context, symbol string) (*exchange.Quote, error)
	SaveQuote(ctx context.Context, quote *exchange.Quote) error
	SubscribeToQuotes(ctx context.Context, symbol string) *redis.PubSub
}

// HistoryRepository defines the interface for historical market data persistence.
type HistoryRepository interface {
	SaveQuote(ctx context.Context, quote *exchange.Quote) error
	GetHistory(ctx context.Context, symbol string, limit int) ([]*exchange.Quote, error)
}

// UserRepository defines the interface for user persistence.
type UserRepository interface {
	GetUsers(ctx context.Context) ([]*user.User, error)
	GetUser(ctx context.Context, id int64) (*user.User, error)
	GetUserByUsername(
		ctx context.Context,
		username string,
	) (*user.User, string, error) // Returns user, hash, error
	CreateUser(
		ctx context.Context,
		username string,
		hashedPassword string,
		firstName string,
		lastName string,
		balance float64,
		website string,
		isPublic bool,
	) (*user.User, error)

	GetUserForUpdate(ctx context.Context, id int64) (*user.User, error)
	UpdateUserProfile(ctx context.Context, user *user.User) error
	UpdateUserBalance(ctx context.Context, id int64, balance float64) error
	WithTx(tx Transaction) UserRepository
}

// PortfolioRepository defines the interface for portfolio persistence.
type PortfolioRepository interface {
	GetPortfolio(ctx context.Context, userID int64) ([]*user.PortfolioItem, error)
	GetPortfolioItem(ctx context.Context, userID int64, symbol string) (*user.PortfolioItem, error)

	GetPortfolioItemForUpdate(
		ctx context.Context,
		userID int64,
		symbol string,
	) (*user.PortfolioItem, error)
	SetPortfolioItem(
		ctx context.Context,
		userID int64,
		symbol string,
		quantity float64,
		averagePrice float64,
	) error
	DeletePortfolioItem(ctx context.Context, userID int64, symbol string) error
	WithTx(tx Transaction) PortfolioRepository
}
