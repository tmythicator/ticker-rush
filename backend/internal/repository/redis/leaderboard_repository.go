// Package redis provides Valkey/Redis repositories.
package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"

	"github.com/tmythicator/ticker-rush/backend/internal/service"
)

const (
	leaderboardPrefix           = "leaderboard"
	leaderboardLastUpdatePrefix = "leaderboard:last_update"
)

func leaderboardKey(ladderID int64) string {
	return fmt.Sprintf("%s:%d", leaderboardPrefix, ladderID)
}

func leaderboardLastUpdateKey(ladderID int64) string {
	return fmt.Sprintf("%s:%d", leaderboardLastUpdatePrefix, ladderID)
}

// LeaderboardRepository handles leaderboard storage in Redis/Valkey.
type LeaderboardRepository struct {
	valkey *redis.Client
}

// NewLeaderboardRepository creates a new instance of LeaderboardRepository.
func NewLeaderboardRepository(valkey *redis.Client) *LeaderboardRepository {
	return &LeaderboardRepository{valkey: valkey}
}

// UpdateRank updates the user's score/rank in the leaderboard.
func (r *LeaderboardRepository) UpdateRank(ctx context.Context, ladderID int64, userID int64, score float64) error {
	key := leaderboardKey(ladderID)
	member := strconv.FormatInt(userID, 10)

	return r.valkey.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: member,
	}).Err()
}

// GetLeaderboard retrieves a range of users and their scores from the leaderboard.
func (r *LeaderboardRepository) GetLeaderboard(ctx context.Context, ladderID int64, offset, limit int) ([]service.LeaderboardScore, error) {
	key := leaderboardKey(ladderID)
	start := int64(offset)
	stop := start + int64(limit) - 1

	res, err := r.valkey.ZRevRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}

	scores := make([]service.LeaderboardScore, len(res))
	for i, item := range res {
		userID, _ := strconv.ParseInt(item.Member.(string), 10, 64)
		scores[i] = service.LeaderboardScore{
			UserID: userID,
			Score:  item.Score,
		}
	}

	return scores, nil
}

// GetLastUpdate retrieves the Unix timestamp of the last leaderboard update.
func (r *LeaderboardRepository) GetLastUpdate(ctx context.Context, ladderID int64) (int64, error) {
	key := leaderboardLastUpdateKey(ladderID)
	val, err := r.valkey.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}

		return 0, err
	}

	timestamp, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return timestamp, nil
}

// SetLastUpdate updates the Unix timestamp of the last leaderboard update.
func (r *LeaderboardRepository) SetLastUpdate(ctx context.Context, ladderID int64, timestamp int64) error {
	key := leaderboardLastUpdateKey(ladderID)

	return r.valkey.Set(ctx, key, strconv.FormatInt(timestamp, 10), 0).Err()
}

// GetTotalCount returns the total number of participants in the leaderboard.
func (r *LeaderboardRepository) GetTotalCount(ctx context.Context, ladderID int64) (int64, error) {
	key := leaderboardKey(ladderID)

	return r.valkey.ZCard(ctx, key).Result()
}

// RemoveFromLeaderboard removes a user from the leaderboard.
func (r *LeaderboardRepository) RemoveFromLeaderboard(ctx context.Context, ladderID int64, userID int64) error {
	key := leaderboardKey(ladderID)
	member := strconv.FormatInt(userID, 10)

	return r.valkey.ZRem(ctx, key, member).Err()
}
