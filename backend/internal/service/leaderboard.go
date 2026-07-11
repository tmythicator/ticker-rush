package service

import (
	"context"
	"time"

	"github.com/tmythicator/ticker-rush/backend/internal/domain"
)

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

// Leaderboard handles the calculation and retrieval of user rankings.
type Leaderboard struct {
	userRepo        UserRepo
	portfolioRepo   PortfolioRepository
	marketRepo      MarketRepository
	leaderboardRepo LeaderboardRepository
	ladderRepo      LadderRepository
}

// NewLeaderboard creates a new instance of Leaderboard with required dependencies.
func NewLeaderboard(
	userRepo UserRepo,
	portfolioRepo PortfolioRepository,
	marketRepo MarketRepository,
	ladderRepo LadderRepository,
	leaderboardRepo LeaderboardRepository,
) *Leaderboard {
	return &Leaderboard{
		userRepo:        userRepo,
		portfolioRepo:   portfolioRepo,
		marketRepo:      marketRepo,
		leaderboardRepo: leaderboardRepo,
		ladderRepo:      ladderRepo,
	}
}

// UpdateLeaderboard recalculates the net worth of all users and updates the Redis Sorted Set for the active ladder.
func (s *Leaderboard) UpdateLeaderboard(ctx context.Context) error {
	ladderID, err := s.ladderRepo.GetActiveLadder(ctx)
	if err != nil {
		return err
	}
	users, err := s.userRepo.GetUsers(ctx)
	if err != nil {
		return err
	}

	for _, u := range users {
		balance, errBalance := s.userRepo.GetUserBalance(ctx, u.ID, ladderID)
		if errBalance != nil {
			continue
		}

		totalWorth := balance

		portfolio, errPortfolio := s.portfolioRepo.GetPortfolio(ctx, u.ID, ladderID)
		if errPortfolio != nil {
			continue
		}

		for _, item := range portfolio {
			quote, errQuote := s.marketRepo.GetQuote(ctx, item.StockSymbol)
			if errQuote != nil {
				continue
			}

			itemValue := quote.Price.Mul(item.Quantity)
			totalWorth = totalWorth.Add(itemValue)
		}

		scoreVal, _ := totalWorth.Float64()
		err = s.leaderboardRepo.UpdateRank(ctx, ladderID, u.ID, scoreVal)
		if err != nil {
			continue
		}
	}

	err = s.leaderboardRepo.SetLastUpdate(ctx, ladderID, time.Now().Unix())
	if err != nil {
		return err
	}

	return nil
}

// GetLeaderboard retrieves a paginated list of top performing users from Redis/Valkey for the active ladder.
func (s *Leaderboard) GetLeaderboard(ctx context.Context, offset, limit int) (*domain.LeaderboardResponse, error) {
	ladderID, err := s.ladderRepo.GetActiveLadder(ctx)
	if err != nil {
		return nil, err
	}

	totalCount, err := s.leaderboardRepo.GetTotalCount(ctx, ladderID)
	if err != nil {
		return nil, err
	}

	lastUpdate, err := s.leaderboardRepo.GetLastUpdate(ctx, ladderID)
	if err != nil {
		lastUpdate = 0
	}

	if totalCount == 0 {
		return &domain.LeaderboardResponse{
			Entries:    []domain.LeaderboardEntry{},
			TotalCount: 0,
			LastUpdate: lastUpdate,
		}, nil
	}

	scores, err := s.leaderboardRepo.GetLeaderboard(ctx, ladderID, offset, limit)
	if err != nil {
		return nil, err
	}

	entries := make([]domain.LeaderboardEntry, 0, len(scores))

	for i, item := range scores {
		fetchedUser, err := s.userRepo.GetUser(ctx, item.UserID)
		if err != nil {
			_ = s.leaderboardRepo.RemoveFromLeaderboard(ctx, ladderID, item.UserID)
			totalCount--

			continue
		}
		leadEntry := domain.LeaderboardEntry{
			User:  s.anonymizeUser(fetchedUser),
			Rank:  int32(i) + int32(offset) + 1,
			Score: item.Score,
		}

		entries = append(entries, leadEntry)
	}

	return &domain.LeaderboardResponse{
		Entries:    entries,
		TotalCount: int32(totalCount),
		LastUpdate: lastUpdate,
	}, nil
}

// anonymizeUser masks user fields if the user is not public.
func (s *Leaderboard) anonymizeUser(u *domain.User) domain.User {
	if u.IsPublic {
		return *u
	}

	return domain.User{
		ID:              u.ID,
		Username:        "Classified",
		FirstName:       "",
		LastName:        "",
		Website:         "",
		IsPublic:        false,
		IsAdmin:         false,
		IsBanned:        u.IsBanned,
		CreatedAt:       u.CreatedAt,
		Balance:         u.Balance,
		Portfolio:       nil,
		IsParticipating: u.IsParticipating,
	}
}
