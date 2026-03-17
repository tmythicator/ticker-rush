package service

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/tmythicator/ticker-rush/backend/internal/proto/leaderboard/v1"
)

// LeaderBoardService handles the calculation and retrieval of user rankings.
type LeaderBoardService struct {
	userRepo      UserRepository
	portfolioRepo PortfolioRepository
	marketRepo    MarketRepository
	redisClient   *redis.Client
	ladderRepo    LadderRepository
}

const leaderboardKeyPrefix = "leaderboard:"
const lastUpdateKeyPrefix = "leaderboard:last_update:"

// NewLeaderBoardService creates a new instance of LeaderBoardService with required dependencies.
func NewLeaderBoardService(
	userRepo UserRepository,
	portfolioRepo PortfolioRepository,
	marketRepo MarketRepository,
	ladderRepo LadderRepository,
	redisClient *redis.Client,
) *LeaderBoardService {
	return &LeaderBoardService{
		userRepo:      userRepo,
		portfolioRepo: portfolioRepo,
		marketRepo:    marketRepo,
		redisClient:   redisClient,
		ladderRepo:    ladderRepo,
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

			totalWorth += quote.Price * item.Quantity
		}

		lKey := leaderboardKeyPrefix + strconv.FormatInt(ladderID, 10)
		err = s.redisClient.ZAdd(ctx, lKey, redis.Z{Score: totalWorth, Member: u.Id}).Err()
		if err != nil {
			continue
		}
	}

	uKey := lastUpdateKeyPrefix + strconv.FormatInt(ladderID, 10)
	err = s.redisClient.Set(ctx, uKey, time.Now().Unix(), 0).Err()
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

	lKey := leaderboardKeyPrefix + strconv.FormatInt(ladderID, 10)
	totalCount, err := s.redisClient.ZCard(ctx, lKey).Result()
	if err != nil {
		return nil, err
	}

	var lastUpdate int64
	uKey := lastUpdateKeyPrefix + strconv.FormatInt(ladderID, 10)
	lastUpdateStr, errTime := s.redisClient.Get(ctx, uKey).Result()
	if errTime == nil {
		lastUpdate, _ = strconv.ParseInt(lastUpdateStr, 10, 64)
	}

	if totalCount == 0 {
		return &leaderboard.GetLeaderboardResponse{
			Entries:    []*leaderboard.LeaderboardEntry{},
			TotalCount: 0,
			LastUpdate: lastUpdate,
		}, nil
	}

	start := int64(req.GetOffset())
	stop := start + int64(req.GetLimit()) - 1

	res, err := s.redisClient.ZRevRangeWithScores(ctx, lKey, start, stop).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]*leaderboard.LeaderboardEntry, 0, len(res))

	for i, item := range res {
		userID, _ := strconv.ParseInt(item.Member.(string), 10, 64)
		fetchedUser, err := s.userRepo.GetUser(ctx, userID)
		if err != nil {
			s.redisClient.ZRem(ctx, lKey, item.Member)
			totalCount--

			continue

		}
		leadEntry := &leaderboard.LeaderboardEntry{
			User:  fetchedUser,
			Rank:  int32(i) + int32(start) + 1,
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
