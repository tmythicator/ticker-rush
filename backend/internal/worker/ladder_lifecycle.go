package worker

import (
	"context"
	"log"
	"sort"
	"time"

	"github.com/shopspring/decimal"

	"github.com/tmythicator/ticker-rush/backend/internal/service"
)

// LadderLifecycleWorker periodically checks for expired and pending ladders to transition them.
type LadderLifecycleWorker struct {
	ladderRepo    service.LadderRepository
	portfolioRepo service.PortfolioRepository
	marketRepo    service.MarketRepository
	interval      time.Duration
}

// NewLadderLifecycleWorker creates a new LadderLifecycleWorker.
func NewLadderLifecycleWorker(
	ladderRepo service.LadderRepository,
	portfolioRepo service.PortfolioRepository,
	marketRepo service.MarketRepository,
	interval time.Duration,
) *LadderLifecycleWorker {
	return &LadderLifecycleWorker{
		ladderRepo:    ladderRepo,
		portfolioRepo: portfolioRepo,
		marketRepo:    marketRepo,
		interval:      interval,
	}
}

// Start runs the periodic check loop.
func (w *LadderLifecycleWorker) Start(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	log.Println("[LadderLifecycleWorker] Performing initial check...")
	w.RunOnce(ctx)

	for {
		select {
		case <-ticker.C:
			w.RunOnce(ctx)
		case <-ctx.Done():
			log.Println("[LadderLifecycleWorker] Stopping worker...")

			return ctx.Err()
		}
	}
}

// RunOnce executes one cycle of deactivation and activation checks.
func (w *LadderLifecycleWorker) RunOnce(ctx context.Context) {
	now := time.Now()
	if err := w.DeactivateExpiredLadders(ctx, now); err != nil {
		log.Printf("[LadderLifecycleWorker] Error deactivating expired ladders: %v", err)
	}
	if err := w.ActivatePendingLadders(ctx, now); err != nil {
		log.Printf("[LadderLifecycleWorker] Error activating pending ladders: %v", err)
	}
}

// DeactivateExpiredLadders finds and deactivates expired ladders, calculating rankings and pruning data.
func (w *LadderLifecycleWorker) DeactivateExpiredLadders(ctx context.Context, now time.Time) error {
	expired, err := w.ladderRepo.GetExpiredActiveLadders(ctx, now)
	if err != nil {
		return err
	}

	for _, l := range expired {
		log.Printf("[LadderLifecycleWorker] Processing expiration of ladder %d (%s)...", l.Id, l.Name)

		// 1. Fetch participants
		participants, err := w.ladderRepo.GetLadderParticipants(ctx, l.Id)
		if err != nil {
			log.Printf("[LadderLifecycleWorker] Failed to get participants for ladder %d: %v", l.Id, err)

			continue
		}

		type participantScore struct {
			userID   int64
			netWorth decimal.Decimal
		}
		var scores []participantScore
		quoteCache := make(map[string]decimal.Decimal)

		for _, p := range participants {
			netWorth := p.Balance // Starts with liquid cash balance

			// Fetch portfolio items
			items, err := w.portfolioRepo.GetPortfolio(ctx, p.UserID, l.Id)
			if err != nil {
				log.Printf("[LadderLifecycleWorker] Failed to get portfolio for user %d in ladder %d: %v", p.UserID, l.Id, err)

				continue
			}

			for _, item := range items {
				symbol := item.GetStockSymbol()
				qty := decimal.NewFromFloat(item.GetQuantity())
				if qty.IsZero() {
					continue
				}

				price, exists := quoteCache[symbol]
				if !exists {
					quote, err := w.marketRepo.GetQuote(ctx, symbol)
					if err != nil {
						log.Printf("[LadderLifecycleWorker] Failed to get quote for %s: %v", symbol, err)
						price = decimal.Zero
					} else {
						price = decimal.NewFromFloat(quote.GetPrice())
					}
					quoteCache[symbol] = price
				}

				netWorth = netWorth.Add(qty.Mul(price))
			}

			scores = append(scores, participantScore{
				userID:   p.UserID,
				netWorth: netWorth,
			})
		}

		// 2. Sort scores in descending order of net worth
		sort.Slice(scores, func(i, j int) bool {
			return scores[i].netWorth.GreaterThan(scores[j].netWorth)
		})

		// 3. Save final scores and ranks
		for rankIdx, score := range scores {
			rank := int32(rankIdx + 1)
			err := w.ladderRepo.InsertLadderParticipant(ctx, l.Id, score.userID, score.netWorth, rank)
			if err != nil {
				log.Printf("[LadderLifecycleWorker] Failed to insert final rank for user %d in ladder %d: %v", score.userID, l.Id, err)
			}
		}

		// 4. Prune portfolio items for this ladder (space saving)
		if err := w.ladderRepo.DeleteLadderPortfolioItemsByLadder(ctx, l.Id); err != nil {
			log.Printf("[LadderLifecycleWorker] Failed to prune portfolio items for ladder %d: %v", l.Id, err)
		}

		// 5. Prune participants outside the top 20
		if err := w.ladderRepo.PruneLadderParticipants(ctx, l.Id, 20); err != nil {
			log.Printf("[LadderLifecycleWorker] Failed to prune participants for ladder %d: %v", l.Id, err)
		}

		// 6. Deactivate ladder
		if err := w.ladderRepo.UpdateLadderStatus(ctx, l.Id, false); err != nil {
			log.Printf("[LadderLifecycleWorker] Failed to deactivate ladder %d: %v", l.Id, err)
		}

		log.Printf("[LadderLifecycleWorker] Ladder %d has been successfully deactivated and pruned.", l.Id)
	}

	return nil
}

// ActivatePendingLadders finds and activates pending ladders.
func (w *LadderLifecycleWorker) ActivatePendingLadders(ctx context.Context, now time.Time) error {
	pending, err := w.ladderRepo.GetPendingLaddersToActivate(ctx, now)
	if err != nil {
		return err
	}

	for _, l := range pending {
		log.Printf("[LadderLifecycleWorker] Activating pending ladder %d (%s)...", l.Id, l.Name)
		if err := w.ladderRepo.UpdateLadderStatus(ctx, l.Id, true); err != nil {
			log.Printf("[LadderLifecycleWorker] Failed to activate ladder %d: %v", l.Id, err)
		} else {
			log.Printf("[LadderLifecycleWorker] Ladder %d is now active.", l.Id)
		}
	}

	return nil
}
