package service

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	"github.com/tmythicator/ticker-rush/backend/internal/domain"
)

// LadderRepository defines the interface for ladder management.
type LadderRepository interface {
	GetActiveLadder(ctx context.Context) (int64, error)
	GetLadder(ctx context.Context, id int64) (*domain.Ladder, error)
	GetAllowedTickers(ctx context.Context, ladderID int64) ([]*domain.TickerInfo, error)
	JoinLadder(ctx context.Context, ladderID int64, userID int64) error
	IsUserInLadder(ctx context.Context, ladderID int64, userID int64) (bool, error)
	GetExpiredActiveLadders(ctx context.Context, now time.Time) ([]*domain.Ladder, error)
	GetPendingLaddersToActivate(ctx context.Context, now time.Time) ([]*domain.Ladder, error)
	UpdateLadderStatus(ctx context.Context, id int64, isActive bool) error
	GetLadderParticipants(ctx context.Context, ladderID int64) ([]domain.LadderParticipant, error)
	InsertLadderParticipant(ctx context.Context, ladderID int64, userID int64, finalBalance decimal.Decimal, finalRank int32) error
	PruneLadderParticipants(ctx context.Context, ladderID int64, rankThreshold int32) error
	DeleteLadderPortfolioItemsByLadder(ctx context.Context, ladderID int64) error
}

// Ladder handles ladder-related business logic.
type Ladder struct {
	ladderRepo LadderRepository
}

// NewLadder creates a new instance of Ladder.
func NewLadder(ladderRepo LadderRepository) *Ladder {
	return &Ladder{
		ladderRepo: ladderRepo,
	}
}

// GetActiveLadderID retrieves the ID of the currently active ladder.
func (s *Ladder) GetActiveLadderID(ctx context.Context) (int64, error) {
	return s.ladderRepo.GetActiveLadder(ctx)
}

// GetAllowedTickers retrieves the allowed stock symbols for the active ladder.
func (s *Ladder) GetAllowedTickers(ctx context.Context) ([]*domain.TickerInfo, error) {
	ladderID, err := s.ladderRepo.GetActiveLadder(ctx)
	if err != nil {
		return nil, err
	}

	return s.ladderRepo.GetAllowedTickers(ctx, ladderID)
}

// GetActiveLadder retrieves full metadata for the currently active ladder.
func (s *Ladder) GetActiveLadder(ctx context.Context) (*domain.Ladder, error) {
	ladderID, err := s.ladderRepo.GetActiveLadder(ctx)
	if err != nil {
		return nil, err
	}

	return s.ladderRepo.GetLadder(ctx, ladderID)
}

// JoinLadder adds the user to the active ladder.
func (s *Ladder) JoinLadder(ctx context.Context, userID int64) error {
	ladderID, err := s.ladderRepo.GetActiveLadder(ctx)
	if err != nil {
		return err
	}

	return s.ladderRepo.JoinLadder(ctx, ladderID, userID)
}
