package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tmythicator/ticker-rush/server/internal/gen/sqlc"
	"github.com/tmythicator/ticker-rush/server/internal/proto/exchange/v1"
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
func (r *HistoryRepository) SaveQuote(ctx context.Context, quote *exchange.Quote) error {
	// Convert float price to numeric
	priceStr := fmt.Sprintf("%f", quote.GetPrice())
	var price pgtype.Numeric
	err := price.Scan(priceStr)
	if err != nil {
		return fmt.Errorf("failed to convert price to numeric: %w", err)
	}

	return r.queries.CreateQuote(ctx, sqlc.CreateQuoteParams{
		Symbol:    quote.GetSymbol(),
		Price:     price,
		Source:    quote.GetSource(),
		CreatedAt: pgtype.Timestamptz{Time: time.Unix(quote.GetTimestamp(), 0).UTC(), Valid: true},
	})
}

// GetHistory retrieves historical quotes for a symbol.
func (r *HistoryRepository) GetHistory(ctx context.Context, symbol string, limit int) ([]*exchange.Quote, error) {
	rows, err := r.queries.GetHistoryForSymbol(ctx, sqlc.GetHistoryForSymbolParams{
		Symbol: symbol,
		Limit:  int32(limit),
	})
	if err != nil {
		return nil, err
	}

	quotes := make([]*exchange.Quote, len(rows))
	for i, row := range rows {
		// Convert numeric price back to float
		priceFloat, _ := row.Price.Float64Value()

		quotes[i] = &exchange.Quote{
			Symbol:    row.Symbol,
			Price:     priceFloat.Float64,
			Source:    row.Source,
			Timestamp: row.CreatedAt.Time.Unix(),
		}
	}

	return quotes, nil
}
