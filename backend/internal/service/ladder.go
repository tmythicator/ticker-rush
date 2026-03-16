package service

import (
	"context"

	"github.com/tmythicator/ticker-rush/backend/internal/proto/ladder/v1"
)

// LadderService handles ladder-related business logic.
type LadderService struct {
	ladderRepo LadderRepository
}

// NewLadderService creates a new instance of LadderService.
func NewLadderService(ladderRepo LadderRepository) *LadderService {
	return &LadderService{
		ladderRepo: ladderRepo,
	}
}

// GetActiveLadderID retrieves the ID of the currently active ladder.
func (s *LadderService) GetActiveLadderID(ctx context.Context) (int64, error) {
	return s.ladderRepo.GetActiveLadder(ctx)
}

// GetAllowedTickers retrieves the allowed stock symbols for the active ladder.
func (s *LadderService) GetAllowedTickers(ctx context.Context) ([]*ladder.TickerInfo, error) {
	ladderID, err := s.ladderRepo.GetActiveLadder(ctx)
	if err != nil {
		return nil, err
	}
	return s.ladderRepo.GetAllowedTickers(ctx, ladderID)
}

// GetActiveLadder retrieves full metadata for the currently active ladder.
func (s *LadderService) GetActiveLadder(ctx context.Context) (*ladder.Ladder, error) {
	ladderID, err := s.ladderRepo.GetActiveLadder(ctx)
	if err != nil {
		return nil, err
	}
	return s.ladderRepo.GetLadder(ctx, ladderID)
}
