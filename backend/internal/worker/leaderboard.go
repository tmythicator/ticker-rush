// Package worker provides background workers.
package worker

import (
	"context"
	"log"
	"time"

	"github.com/tmythicator/ticker-rush/server/internal/service"
)

// LeaderboardWorker is a worker that periodically updates the leaderboard.
type LeaderboardWorker struct {
	lbService *service.LeaderBoardService
	interval  time.Duration
}

// NewLeaderboardWorker creates a new instance of LeaderboardWorker.
func NewLeaderboardWorker(lbService *service.LeaderBoardService, interval time.Duration) *LeaderboardWorker {
	return &LeaderboardWorker{
		lbService: lbService,
		interval:  interval,
	}
}

// Start begins the leaderboard update loop.
func (w *LeaderboardWorker) Start(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	log.Println("[LeaderboardWorker] Performing initial update...")
	if err := w.lbService.UpdateLeaderboard(ctx); err != nil {
		log.Printf("[LeaderboardWorker] Initial update failed: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			log.Println("[LeaderboardWorker] Updating leaderboard...")
			if err := w.lbService.UpdateLeaderboard(ctx); err != nil {
				log.Printf("[LeaderboardWorker] Update failed: %v", err)
			}
		case <-ctx.Done():
			log.Println("[LeaderboardWorker] Stopping...")

			return ctx.Err()
		}
	}
}
