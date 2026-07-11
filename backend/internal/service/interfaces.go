package service

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	"github.com/tmythicator/ticker-rush/backend/internal/domain"
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

// CreateUserParams represents parameters for creating a new user.
type CreateUserParams struct {
	Username      string
	PasswordHash  string
	FirstName     string
	LastName      string
	Website       string
	IsPublic      bool
	AgbAcceptedAt time.Time
}

// PortfolioRepository defines the interface for portfolio persistence.
type PortfolioRepository interface {
	GetPortfolio(ctx context.Context, userID int64, ladderID int64) ([]*domain.PortfolioItem, error)
	GetPortfolioItem(ctx context.Context, userID int64, ladderID int64, symbol string) (*domain.PortfolioItem, error)

	GetPortfolioItemForUpdate(
		ctx context.Context,
		userID int64,
		ladderID int64,
		symbol string,
	) (*domain.PortfolioItem, error)
	SetPortfolioItem(
		ctx context.Context,
		userID int64,
		ladderID int64,
		symbol string,
		quantity decimal.Decimal,
		averagePrice decimal.Decimal,
	) error
	DeletePortfolioItem(ctx context.Context, userID int64, ladderID int64, symbol string) error
	WithTx(tx Transaction) PortfolioRepository
}

// LeaderboardScore represents a single score entry in the leaderboard.
type LeaderboardScore struct {
	UserID int64
	Score  float64
}

// LeaderboardRepository defines the interface for leaderboard state storage (e.g. Valkey/Redis).
type LeaderboardRepository interface {
	UpdateRank(ctx context.Context, ladderID int64, userID int64, score float64) error
	GetLeaderboard(ctx context.Context, ladderID int64, offset, limit int) ([]LeaderboardScore, error)
	GetLastUpdate(ctx context.Context, ladderID int64) (int64, error)
	SetLastUpdate(ctx context.Context, ladderID int64, timestamp int64) error
	GetTotalCount(ctx context.Context, ladderID int64) (int64, error)
	RemoveFromLeaderboard(ctx context.Context, ladderID int64, userID int64) error
}
