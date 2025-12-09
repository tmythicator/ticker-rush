package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/tmythicator/ticker-rush/server/internal/gen/sqlc"
)

type PortfolioRepository struct {
	queries *db.Queries
}

func (r *PortfolioRepository) WithTx(tx pgx.Tx) *PortfolioRepository {
	return &PortfolioRepository{
		queries: r.queries.WithTx(tx),
	}
}

func NewPortfolioRepository(pool *pgxpool.Pool) *PortfolioRepository {
	return &PortfolioRepository{
		queries: db.New(pool),
	}
}

func (r *PortfolioRepository) GetPortfolio(ctx context.Context, userID int64) ([]db.PortfolioItem, error) {
	return r.queries.GetPortfolio(ctx, userID)
}

func (r *PortfolioRepository) GetPortfolioItem(ctx context.Context, userID int64, symbol string) (db.PortfolioItem, error) {
	return r.queries.GetPortfolioItem(ctx, db.GetPortfolioItemParams{
		UserID:      userID,
		StockSymbol: symbol,
	})
}

func (r *PortfolioRepository) GetPortfolioItemForUpdate(ctx context.Context, userID int64, symbol string) (db.PortfolioItem, error) {
	return r.queries.GetPortfolioItemForUpdate(ctx, db.GetPortfolioItemForUpdateParams{
		UserID:      userID,
		StockSymbol: symbol,
	})
}

func (r *PortfolioRepository) SetPortfolioItem(ctx context.Context, userID int64, symbol string, quantity float64, averagePrice float64) error {
	return r.queries.SetPortfolioItem(ctx, db.SetPortfolioItemParams{
		UserID:       userID,
		StockSymbol:  symbol,
		Quantity:     quantity,
		AveragePrice: averagePrice,
	})
}

func (r *PortfolioRepository) DeletePortfolioItem(ctx context.Context, userID int64, symbol string) error {
	return r.queries.DeletePortfolioItem(ctx, db.DeletePortfolioItemParams{
		UserID:      userID,
		StockSymbol: symbol,
	})
}
