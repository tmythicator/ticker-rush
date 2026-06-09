package service

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"

	"github.com/tmythicator/ticker-rush/backend/internal/gen/sqlc"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/exchange/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/ladder/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/portfolio/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/user/v1"
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

// UserRepository defines the interface for user persistence.
type UserRepository interface {
	GetUsers(ctx context.Context) ([]*user.User, error)
	GetUser(ctx context.Context, id int64) (*user.User, error)
	GetUserByUsername(
		ctx context.Context,
		username string,
	) (*user.User, string, error) // Returns user, hash, error
	CreateUser(ctx context.Context, params CreateUserParams) (*user.User, error)

	GetUserForUpdate(ctx context.Context, id int64) (*user.User, error)
	UpdateUserProfile(ctx context.Context, user *user.User) error
	UpdateUserBalance(ctx context.Context, userID int64, ladderID int64, balance decimal.Decimal) error
	GetUserBalance(ctx context.Context, userID int64, ladderID int64) (decimal.Decimal, error)
	GetUserWithPortfolioForActiveLadder(ctx context.Context, id int64) ([]sqlc.GetUserWithPortfolioForActiveLadderRow, error)
	WithTx(tx Transaction) UserRepository
}

// LadderRepository defines the interface for ladder management.
type LadderRepository interface {
	GetActiveLadder(ctx context.Context) (int64, error)
	GetLadder(ctx context.Context, id int64) (*ladder.Ladder, error)
	GetAllowedTickers(ctx context.Context, ladderID int64) ([]*ladder.TickerInfo, error)
	JoinLadder(ctx context.Context, ladderID int64, userID int64) error
	IsUserInLadder(ctx context.Context, ladderID int64, userID int64) (bool, error)
	GetExpiredActiveLadders(ctx context.Context, now time.Time) ([]*ladder.Ladder, error)
	GetPendingLaddersToActivate(ctx context.Context, now time.Time) ([]*ladder.Ladder, error)
	UpdateLadderStatus(ctx context.Context, id int64, isActive bool) error
	GetLadderParticipants(ctx context.Context, ladderID int64) ([]sqlc.LadderParticipant, error)
	InsertLadderParticipant(ctx context.Context, ladderID int64, userID int64, finalBalance decimal.Decimal, finalRank int32) error
	PruneLadderParticipants(ctx context.Context, ladderID int64, rankThreshold int32) error
	DeleteLadderPortfolioItemsByLadder(ctx context.Context, ladderID int64) error
	WithTx(tx Transaction) LadderRepository
}

// PortfolioRepository defines the interface for portfolio persistence.
type PortfolioRepository interface {
	GetPortfolio(ctx context.Context, userID int64, ladderID int64) ([]*portfolio.PortfolioItem, error)
	GetPortfolioItem(ctx context.Context, userID int64, ladderID int64, symbol string) (*portfolio.PortfolioItem, error)

	GetPortfolioItemForUpdate(
		ctx context.Context,
		userID int64,
		ladderID int64,
		symbol string,
	) (*portfolio.PortfolioItem, error)
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
