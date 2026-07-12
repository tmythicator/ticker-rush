// Package domain contains domain models and business types.
package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// User represents a system user and their active participation state.
type User struct {
	ID              int64
	Username        string
	FirstName       string
	LastName        string
	Website         string
	IsPublic        bool
	IsAdmin         bool
	IsBanned        bool
	CreatedAt       time.Time
	Balance         decimal.Decimal
	Portfolio       map[string]PortfolioItem
	IsParticipating bool
}

// PortfolioItem represents stock holdings of a user.
type PortfolioItem struct {
	StockSymbol  string
	Quantity     decimal.Decimal
	AveragePrice decimal.Decimal
}

// TickerInfo represents ticker symbol configurations allowed in ladders.
type TickerInfo struct {
	Symbol string
	Source string
}

// Ladder represents a competition cycle.
type Ladder struct {
	ID             int64
	Name           string
	Type           string
	StartTime      time.Time
	EndTime        time.Time
	IsActive       bool
	CreatedAt      time.Time
	InitialBalance decimal.Decimal
	AllowedTickers []TickerInfo
}

// LadderParticipant represents a user's standing in a ladder.
type LadderParticipant struct {
	LadderID     int64
	User         User
	Balance      decimal.Decimal
	FinalRank    int32
	FinalBalance decimal.Decimal
	JoinedAt     time.Time
}

// Quote represents a ticker quote.
type Quote struct {
	Symbol        string
	Price         decimal.Decimal
	Change        decimal.Decimal
	ChangePercent decimal.Decimal
	Timestamp     time.Time
	Source        string
	IsClosed      bool
}

// LeaderboardEntry represents a single rank entry on the leaderboard.
type LeaderboardEntry struct {
	User  User
	Rank  int32
	Score float64
}

// LeaderboardResponse represents a paginated leaderboard response.
type LeaderboardResponse struct {
	Entries    []LeaderboardEntry
	TotalCount int32
	LastUpdate int64
}

// MarketStalenessThreshold defines the maximum time allowed since the last update
// before the market data is considered closed or stale.
const MarketStalenessThreshold = 30 * time.Minute

// CalculateIsClosed determines if a market is closed based on a timestamp and previous status.
func CalculateIsClosed(originalClosed bool, timestamp time.Time) bool {
	return originalClosed || time.Since(timestamp) > MarketStalenessThreshold
}
