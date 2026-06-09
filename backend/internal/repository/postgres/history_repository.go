package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tmythicator/ticker-rush/backend/internal/domain"
	"github.com/tmythicator/ticker-rush/backend/internal/gen/sqlc"
)

// HistoryRepository implements service.HistoryRepository for PostgreSQL.
type HistoryRepository struct {
	queries *sqlc.Queries
}

// NewHistoryRepository creates a new PostgreSQL HistoryRepository.
func NewHistoryRepository(pool *pgxpool.Pool) *HistoryRepository {
	return &HistoryRepository{
		queries: sqlc.New(pool),
	}
}

// SaveQuote saves a quote to the history table.
func (r *HistoryRepository) SaveQuote(ctx context.Context, quote *domain.Quote) error {
	return r.queries.CreateQuote(ctx, sqlc.CreateQuoteParams{
		Symbol:    quote.Symbol,
		Price:     quote.Price,
		Source:    quote.Source,
		CreatedAt: pgtype.Timestamptz{Time: quote.Timestamp, Valid: true},
	})
}

// GetHistory retrieves historical quotes for a symbol.
func (r *HistoryRepository) GetHistory(ctx context.Context, symbol string, limit int) ([]*domain.Quote, error) {
	rows, err := r.queries.GetHistoryForSymbol(ctx, sqlc.GetHistoryForSymbolParams{
		Symbol: symbol,
		Limit:  int32(limit),
	})
	if err != nil {
		return nil, err
	}

	quotes := make([]*domain.Quote, len(rows))
	for i, row := range rows {
		quotes[i] = &domain.Quote{
			Symbol:    row.Symbol,
			Price:     row.Price,
			Source:    row.Source,
			Timestamp: row.CreatedAt.Time,
		}
	}

	return quotes, nil
}
