package service

import (
	"context"

	"github.com/tmythicator/ticker-rush/backend/internal/apperrors"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/ladder/v1"
)

// LadderService handles ladder-related business logic.
type LadderService struct {
	ladderRepo LadderRepository
	transactor Transactor
}

// NewLadderService creates a new instance of LadderService.
func NewLadderService(ladderRepo LadderRepository, transactor Transactor) *LadderService {
	return &LadderService{
		ladderRepo: ladderRepo,
		transactor: transactor,
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

// JoinLadder adds the user to the active ladder.
func (s *LadderService) JoinLadder(ctx context.Context, userID int64) error {
	ladderID, err := s.ladderRepo.GetActiveLadder(ctx)
	if err != nil {
		return err
	}

	joined, err := s.ladderRepo.IsUserInLadder(ctx, ladderID, userID)
	if err != nil {
		return err
	}
	if joined {
		return apperrors.ErrAlreadyJoinedLadder
	}

	tx, err := s.transactor.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	txLadderRepo := s.ladderRepo.WithTx(tx)

	if err := txLadderRepo.JoinLadder(ctx, ladderID, userID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
