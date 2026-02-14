package service

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/tmythicator/ticker-rush/server/internal/proto/leaderboard/v1"
)

// LeaderBoardService handles the calculation and retrieval of user rankings.
type LeaderBoardService struct {
	userRepo      UserRepository
	portfolioRepo PortfolioRepository
	marketRepo    MarketRepository
	redisClient   *redis.Client
}

const leaderboardKey = "leaderboard"

// NewLeaderBoardService creates a new instance of LeaderBoardService with required dependencies.
func NewLeaderBoardService(
	userRepo UserRepository,
	portfolioRepo PortfolioRepository,
	marketRepo MarketRepository,
	redisClient *redis.Client,
) *LeaderBoardService {
	return &LeaderBoardService{
		userRepo:      userRepo,
		portfolioRepo: portfolioRepo,
		marketRepo:    marketRepo,
		redisClient:   redisClient,
	}
}

// UpdateLeaderboard recalculates the net worth of all users and updates the Redis Sorted Set.
// It iterates through all users, calculates their portfolio value using real-time market quotes,
// and saves the final score to Redis.
func (s *LeaderBoardService) UpdateLeaderboard(ctx context.Context) error {
	users, err := s.userRepo.GetUsers(ctx)
	if err != nil {
		return err
	}

	for _, u := range users {

		totalWorth := u.Balance

		portfolio, err := s.portfolioRepo.GetPortfolio(ctx, u.Id)
		if err != nil {
			continue
		}

		for _, item := range portfolio {
			quote, errQuote := s.marketRepo.GetQuote(ctx, item.StockSymbol)
			if errQuote != nil {
				continue
			}

			totalWorth += quote.Price * item.Quantity
		}

		err = s.redisClient.ZAdd(ctx, leaderboardKey, redis.Z{Score: totalWorth, Member: u.Id}).Err()
		if err != nil {
			continue
		}
	}

	return nil
}

// GetLeaderboard retrieves a paginated list of top performing users from Redis/Valkey.
// It fetches the user details from the database for each entry and handles stale cache cleanup
// if a user is no longer present in the system.
func (s *LeaderBoardService) GetLeaderboard(ctx context.Context, req *leaderboard.GetLeaderboardRequest) (*leaderboard.GetLeaderboardResponse, error) {

	totalCount, err := s.redisClient.ZCard(ctx, leaderboardKey).Result()
	if err != nil {
		return nil, err
	}

	start := int64(req.GetOffset())
	stop := start + int64(req.GetLimit()) - 1

	res, err := s.redisClient.ZRevRangeWithScores(ctx, leaderboardKey, start, stop).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]*leaderboard.LeaderboardEntry, 0, len(res))

	for i, item := range res {
		userID, _ := strconv.ParseInt(item.Member.(string), 10, 64)
		user, err := s.userRepo.GetUser(ctx, userID)
		if err != nil {
			s.redisClient.ZRem(ctx, leaderboardKey, item.Member)
			totalCount--

			continue
		}
		leadEntry := &leaderboard.LeaderboardEntry{
			UserId:        user.Id,
			Rank:          int32(i) + int32(start) + 1,
			TotalNetWorth: item.Score,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
		}
		entries = append(entries, leadEntry)
	}

	return &leaderboard.GetLeaderboardResponse{Entries: entries, TotalCount: int32(totalCount)}, nil
}
