package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/tmythicator/ticker-rush/server/internal/gen/sqlc"
)

type StockRepository struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func NewStockRepository(pool *pgxpool.Pool) *StockRepository {
	return &StockRepository{
		queries: db.New(pool),
		pool:    pool,
	}
}

func (r *StockRepository) UpsertStock(ctx context.Context, symbol string, name string) error {
	_, err := r.queries.UpsertStock(ctx, db.UpsertStockParams{
		Symbol: symbol,
		Name:   name,
	})
	return err
}
