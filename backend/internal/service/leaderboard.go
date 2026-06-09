package service

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	"github.com/tmythicator/ticker-rush/backend/internal/proto/leaderboard/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/user/v1"
)

// LeaderBoardService handles the calculation and retrieval of user rankings.
type LeaderBoardService struct {
	userRepo        UserRepository
	portfolioRepo   PortfolioRepository
	marketRepo      MarketRepository
	leaderboardRepo LeaderboardRepository
	ladderRepo      LadderRepository
}

// NewLeaderBoardService creates a new instance of LeaderBoardService with required dependencies.
func NewLeaderBoardService(
	userRepo UserRepository,
	portfolioRepo PortfolioRepository,
	marketRepo MarketRepository,
	ladderRepo LadderRepository,
	leaderboardRepo LeaderboardRepository,
) *LeaderBoardService {
	return &LeaderBoardService{
		userRepo:        userRepo,
		portfolioRepo:   portfolioRepo,
		marketRepo:      marketRepo,
		leaderboardRepo: leaderboardRepo,
		ladderRepo:      ladderRepo,
	}
}

// UpdateLeaderboard recalculates the net worth of all users and updates the Redis Sorted Set for the active ladder.
// It iterates through all users, calculates their portfolio value using real-time market quotes,
// and saves the final score to Redis.
func (s *LeaderBoardService) UpdateLeaderboard(ctx context.Context) error {
	ladderID, err := s.ladderRepo.GetActiveLadder(ctx)
	if err != nil {
		return err
	}
	users, err := s.userRepo.GetUsers(ctx)
	if err != nil {
		return err
	}

	for _, u := range users {
		balance, errBalance := s.userRepo.GetUserBalance(ctx, u.Id, ladderID)
		if errBalance != nil {
			continue
		}

		totalWorth := balance

		portfolio, errPortfolio := s.portfolioRepo.GetPortfolio(ctx, u.Id, ladderID)
		if errPortfolio != nil {
			continue
		}

		for _, item := range portfolio {
			quote, errQuote := s.marketRepo.GetQuote(ctx, item.StockSymbol)
			if errQuote != nil {
				continue
			}

			itemValue := decimal.NewFromFloat(quote.Price).Mul(decimal.NewFromFloat(item.Quantity))
			totalWorth = totalWorth.Add(itemValue)
		}

		scoreVal, _ := totalWorth.Float64()
		err = s.leaderboardRepo.UpdateRank(ctx, ladderID, u.Id, scoreVal)
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
// It fetches the user details from the database for each entry and handles stale cache cleanup
// if a user is no longer present in the system.
func (s *LeaderBoardService) GetLeaderboard(ctx context.Context, req *leaderboard.GetLeaderboardRequest) (*leaderboard.GetLeaderboardResponse, error) {
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
		return &leaderboard.GetLeaderboardResponse{
			Entries:    []*leaderboard.LeaderboardEntry{},
			TotalCount: 0,
			LastUpdate: lastUpdate,
		}, nil
	}

	scores, err := s.leaderboardRepo.GetLeaderboard(ctx, ladderID, int(req.GetOffset()), int(req.GetLimit()))
	if err != nil {
		return nil, err
	}

	entries := make([]*leaderboard.LeaderboardEntry, 0, len(scores))

	for i, item := range scores {
		fetchedUser, err := s.userRepo.GetUser(ctx, item.UserID)
		if err != nil {
			s.leaderboardRepo.RemoveFromLeaderboard(ctx, ladderID, item.UserID)
			totalCount--

			continue
		}
		leadEntry := &leaderboard.LeaderboardEntry{
			User:  s.anonymizeUser(fetchedUser),
			Rank:  int32(i) + req.GetOffset() + 1,
			Score: item.Score,
		}

		entries = append(entries, leadEntry)
	}

	return &leaderboard.GetLeaderboardResponse{
		Entries:    entries,
		TotalCount: int32(totalCount),
		LastUpdate: lastUpdate,
	}, nil
}

// anonymizeUser masks user fields if the user is not public.
func (s *LeaderBoardService) anonymizeUser(u *user.User) *user.User {
	if u.IsPublic {
		return u
	}

	return &user.User{
		Id:              u.Id,
		Username:        "Classified",
		FirstName:       "",
		LastName:        "",
		Website:         "",
		IsPublic:        false,
		IsAdmin:         false,
		IsBanned:        u.IsBanned,
		CreatedAt:       u.CreatedAt,
		Balance:         0,
		Portfolio:       nil,
		IsParticipating: u.IsParticipating,
	}
}
