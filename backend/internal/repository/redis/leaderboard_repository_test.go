package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	redisRepo "github.com/tmythicator/ticker-rush/backend/internal/repository/redis"
)

func TestLeaderboardRepository(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	rClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer rClient.Close()

	repo := redisRepo.NewLeaderboardRepository(rClient)
	ctx := context.Background()

	const (
		ladderID int64 = 100
		user1    int64 = 1
		user2    int64 = 2
		user3    int64 = 3
	)

	// 1. Test UpdateRank & GetLeaderboard
	err = repo.UpdateRank(ctx, ladderID, user1, 1500.50)
	assert.NoError(t, err)

	err = repo.UpdateRank(ctx, ladderID, user2, 2500.75)
	assert.NoError(t, err)

	err = repo.UpdateRank(ctx, ladderID, user3, 1000.00)
	assert.NoError(t, err)

	// Leaderboard should be ordered descending: user2 (2500.75), user1 (1500.50), user3 (1000.00)
	board, err := repo.GetLeaderboard(ctx, ladderID, 0, 10)
	assert.NoError(t, err)
	assert.Len(t, board, 3)

	assert.Equal(t, user2, board[0].UserID)
	assert.Equal(t, 2500.75, board[0].Score)

	assert.Equal(t, user1, board[1].UserID)
	assert.Equal(t, 1500.50, board[1].Score)

	assert.Equal(t, user3, board[2].UserID)
	assert.Equal(t, 1000.00, board[2].Score)

	// 2. Test GetTotalCount
	count, err := repo.GetTotalCount(ctx, ladderID)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)

	// 3. Test LastUpdate methods
	lastUpd, err := repo.GetLastUpdate(ctx, ladderID)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), lastUpd)

	now := time.Now().Unix()
	err = repo.SetLastUpdate(ctx, ladderID, now)
	assert.NoError(t, err)

	lastUpd, err = repo.GetLastUpdate(ctx, ladderID)
	assert.NoError(t, err)
	assert.Equal(t, now, lastUpd)

	// 4. Test RemoveFromLeaderboard
	err = repo.RemoveFromLeaderboard(ctx, ladderID, user1)
	assert.NoError(t, err)

	count, err = repo.GetTotalCount(ctx, ladderID)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)

	board, err = repo.GetLeaderboard(ctx, ladderID, 0, 10)
	assert.NoError(t, err)
	assert.Len(t, board, 2)
	assert.Equal(t, user2, board[0].UserID)
	assert.Equal(t, user3, board[1].UserID)
}
