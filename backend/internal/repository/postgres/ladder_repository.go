package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/tmythicator/ticker-rush/backend/internal/gen/sqlc"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/ladder/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/service"
)

// LadderRepository handles ladder data persistence.
type LadderRepository struct {
	queries *sqlc.Queries
}

// NewLadderRepository creates a new instance of LadderRepository.
func NewLadderRepository(pool *pgxpool.Pool) *LadderRepository {
	return &LadderRepository{
		queries: sqlc.New(pool),
	}
}

// GetActiveLadder retrieves the ID of the currently active ladder.
func (r *LadderRepository) GetActiveLadder(ctx context.Context) (int64, error) {
	ladder, err := r.queries.GetActiveLadder(ctx)
	if err != nil {
		return 0, err
	}

	return ladder.ID, nil
}

// GetLadder retrieves a ladder by ID.
func (r *LadderRepository) GetLadder(ctx context.Context, id int64) (*ladder.Ladder, error) {
	row, err := r.queries.GetLadder(ctx, id)
	if err != nil {
		return nil, err
	}

	tickers, _ := r.GetAllowedTickers(ctx, id)

	return &ladder.Ladder{
		Id:             row.ID,
		Name:           row.Name,
		Type:           row.Type,
		StartTime:      &timestamppb.Timestamp{Seconds: row.StartTime.Time.Unix()},
		EndTime:        &timestamppb.Timestamp{Seconds: row.EndTime.Time.Unix()},
		IsActive:       row.IsActive,
		InitialBalance: row.InitialBalance.InexactFloat64(),
		AllowedTickers: tickers,
		CreatedAt:      &timestamppb.Timestamp{Seconds: row.CreatedAt.Time.Unix()},
	}, nil
}

// GetAllowedTickers retrieves the allowed stock symbols for a given ladder.
func (r *LadderRepository) GetAllowedTickers(ctx context.Context, ladderID int64) ([]*ladder.TickerInfo, error) {
	tickers, err := r.queries.GetLadderTickers(ctx, ladderID)
	if err != nil {
		return nil, err
	}
	tickerInfos := make([]*ladder.TickerInfo, len(tickers))
	for i, t := range tickers {
		tickerInfos[i] = &ladder.TickerInfo{
			Symbol: t.StockSymbol,
			Source: t.Source,
		}
	}

	return tickerInfos, nil
}

// JoinLadder adds a user to a ladder and initializes their balance.
func (r *LadderRepository) JoinLadder(ctx context.Context, ladderID int64, userID int64) error {
	return r.queries.JoinLadderParticipant(ctx, sqlc.JoinLadderParticipantParams{
		LadderID: ladderID,
		UserID:   userID,
	})
}

// IsUserInLadder checks if a user is enrolled in a ladder.
func (r *LadderRepository) IsUserInLadder(ctx context.Context, ladderID int64, userID int64) (bool, error) {
	return r.queries.IsUserInLadder(ctx, sqlc.IsUserInLadderParams{
		LadderID: ladderID,
		UserID:   userID,
	})
}

// GetExpiredActiveLadders retrieves active ladders whose end time is in the past.
func (r *LadderRepository) GetExpiredActiveLadders(ctx context.Context, now time.Time) ([]*ladder.Ladder, error) {
	rows, err := r.queries.GetExpiredActiveLadders(ctx, pgtype.Timestamptz{Time: now, Valid: true})
	if err != nil {
		return nil, err
	}
	res := make([]*ladder.Ladder, len(rows))
	for i, row := range rows {
		res[i] = &ladder.Ladder{
			Id:             row.ID,
			Name:           row.Name,
			Type:           row.Type,
			StartTime:      &timestamppb.Timestamp{Seconds: row.StartTime.Time.Unix()},
			EndTime:        &timestamppb.Timestamp{Seconds: row.EndTime.Time.Unix()},
			IsActive:       row.IsActive,
			InitialBalance: row.InitialBalance.InexactFloat64(),
			CreatedAt:      &timestamppb.Timestamp{Seconds: row.CreatedAt.Time.Unix()},
		}
	}
	return res, nil
}

// GetPendingLaddersToActivate retrieves inactive ladders whose start time has arrived.
func (r *LadderRepository) GetPendingLaddersToActivate(ctx context.Context, now time.Time) ([]*ladder.Ladder, error) {
	rows, err := r.queries.GetPendingLaddersToActivate(ctx, pgtype.Timestamptz{Time: now, Valid: true})
	if err != nil {
		return nil, err
	}
	res := make([]*ladder.Ladder, len(rows))
	for i, row := range rows {
		res[i] = &ladder.Ladder{
			Id:             row.ID,
			Name:           row.Name,
			Type:           row.Type,
			StartTime:      &timestamppb.Timestamp{Seconds: row.StartTime.Time.Unix()},
			EndTime:        &timestamppb.Timestamp{Seconds: row.EndTime.Time.Unix()},
			IsActive:       row.IsActive,
			InitialBalance: row.InitialBalance.InexactFloat64(),
			CreatedAt:      &timestamppb.Timestamp{Seconds: row.CreatedAt.Time.Unix()},
		}
	}
	return res, nil
}

// UpdateLadderStatus updates the is_active status of a ladder.
func (r *LadderRepository) UpdateLadderStatus(ctx context.Context, id int64, isActive bool) error {
	return r.queries.UpdateLadderStatus(ctx, sqlc.UpdateLadderStatusParams{
		ID:       id,
		IsActive: isActive,
	})
}

// GetLadderParticipants retrieves the participants of a ladder.
func (r *LadderRepository) GetLadderParticipants(ctx context.Context, ladderID int64) ([]sqlc.LadderParticipant, error) {
	return r.queries.GetLadderParticipants(ctx, ladderID)
}

// InsertLadderParticipant updates or inserts a participant's final stats.
func (r *LadderRepository) InsertLadderParticipant(ctx context.Context, ladderID int64, userID int64, finalBalance decimal.Decimal, finalRank int32) error {
	return r.queries.InsertLadderParticipant(ctx, sqlc.InsertLadderParticipantParams{
		LadderID:     ladderID,
		UserID:       userID,
		FinalBalance: finalBalance,
		FinalRank:    pgtype.Int4{Int32: finalRank, Valid: true},
	})
}

// PruneLadderParticipants deletes participants whose final rank is worse than the threshold.
func (r *LadderRepository) PruneLadderParticipants(ctx context.Context, ladderID int64, rankThreshold int32) error {
	return r.queries.PruneLadderParticipants(ctx, sqlc.PruneLadderParticipantsParams{
		LadderID:  ladderID,
		FinalRank: pgtype.Int4{Int32: rankThreshold, Valid: true},
	})
}

// DeleteLadderPortfolioItemsByLadder deletes all stock holdings associated with a ladder.
func (r *LadderRepository) DeleteLadderPortfolioItemsByLadder(ctx context.Context, ladderID int64) error {
	return r.queries.DeleteLadderPortfolioItemsByLadder(ctx, ladderID)
}

// WithTx returns a new LadderRepository that uses the given transaction.
func (r *LadderRepository) WithTx(tx service.Transaction) service.LadderRepository {
	return &LadderRepository{
		queries: r.queries.WithTx(tx.(pgx.Tx)),
	}
}
