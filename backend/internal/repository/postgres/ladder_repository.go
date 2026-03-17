package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
		InitialBalance: row.InitialBalance,
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
	err := r.queries.JoinLadderParticipant(ctx, sqlc.JoinLadderParticipantParams{
		LadderID: ladderID,
		UserID:   userID,
	})
	if err != nil {
		return err
	}

	l, err := r.queries.GetLadder(ctx, ladderID)
	if err != nil {
		return err
	}

	err = r.queries.InsertLadderBalance(ctx, sqlc.InsertLadderBalanceParams{
		LadderID: ladderID,
		UserID:   userID,
		Balance:  l.InitialBalance,
	})
	if err != nil {
		return err
	}

	return nil
}

// IsUserInLadder checks if a user is enrolled in a ladder.
func (r *LadderRepository) IsUserInLadder(ctx context.Context, ladderID int64, userID int64) (bool, error) {
	return r.queries.IsUserInLadder(ctx, sqlc.IsUserInLadderParams{
		LadderID: ladderID,
		UserID:   userID,
	})
}

// WithTx returns a new LadderRepository that uses the given transaction.
func (r *LadderRepository) WithTx(tx service.Transaction) service.LadderRepository {
	return &LadderRepository{
		queries: r.queries.WithTx(tx.(pgx.Tx)),
	}
}
