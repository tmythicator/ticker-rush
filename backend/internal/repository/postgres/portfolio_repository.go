// Package postgres provides PostgreSQL repositories.
package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tmythicator/ticker-rush/backend/internal/gen/sqlc"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/ladder/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/service"
)

// PortfolioRepository handles portfolio data persistence in PostgreSQL.
type PortfolioRepository struct {
	queries *sqlc.Queries
}

// NewPortfolioRepository creates a new instance of PortfolioRepository.
func NewPortfolioRepository(pool *pgxpool.Pool) *PortfolioRepository {
	return &PortfolioRepository{
		queries: sqlc.New(pool),
	}
}

// WithTx returns a new PortfolioRepository that uses the given transaction.
func (r *PortfolioRepository) WithTx(tx service.Transaction) service.PortfolioRepository {
	return &PortfolioRepository{
		queries: r.queries.WithTx(tx.(pgx.Tx)),
	}
}

// GetPortfolio retrieves the portfolio for a user in a specific ladder.
func (r *PortfolioRepository) GetPortfolio(
	ctx context.Context,
	ladderID int64,
	userID int64,
) ([]*ladder.PortfolioItem, error) {
	items, err := r.queries.GetPortfolio(ctx, sqlc.GetPortfolioParams{
		LadderID: ladderID,
		UserID:   userID,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*ladder.PortfolioItem, len(items))
	for i, item := range items {
		result[i] = &ladder.PortfolioItem{
			StockSymbol:  item.StockSymbol,
			Quantity:     item.Quantity,
			AveragePrice: item.AveragePrice,
		}
	}

	return result, nil
}

// GetPortfolioItem retrieves a specific portfolio item.
func (r *PortfolioRepository) GetPortfolioItem(
	ctx context.Context,
	ladderID int64,
	userID int64,
	symbol string,
) (*ladder.PortfolioItem, error) {
	item, err := r.queries.GetPortfolioItem(ctx, sqlc.GetPortfolioItemParams{
		LadderID:    ladderID,
		UserID:      userID,
		StockSymbol: symbol,
	})
	if err != nil {
		return nil, err
	}

	return &ladder.PortfolioItem{
		StockSymbol:  item.StockSymbol,
		Quantity:     item.Quantity,
		AveragePrice: item.AveragePrice,
	}, nil
}

// GetPortfolioItemForUpdate retrieves a portfolio item with a lock for update.
func (r *PortfolioRepository) GetPortfolioItemForUpdate(
	ctx context.Context,
	ladderID int64,
	userID int64,
	symbol string,
) (*ladder.PortfolioItem, error) {
	item, err := r.queries.GetPortfolioItemForUpdate(ctx, sqlc.GetPortfolioItemForUpdateParams{
		LadderID:    ladderID,
		UserID:      userID,
		StockSymbol: symbol,
	})
	if err != nil {
		return nil, err
	}

	return &ladder.PortfolioItem{
		StockSymbol:  item.StockSymbol,
		Quantity:     item.Quantity,
		AveragePrice: item.AveragePrice,
	}, nil
}

// SetPortfolioItem updates or inserts a portfolio item.
func (r *PortfolioRepository) SetPortfolioItem(
	ctx context.Context,
	ladderID int64,
	userID int64,
	symbol string,
	quantity float64,
	averagePrice float64,
) error {
	return r.queries.SetPortfolioItem(ctx, sqlc.SetPortfolioItemParams{
		LadderID:     ladderID,
		UserID:       userID,
		StockSymbol:  symbol,
		Quantity:     quantity,
		AveragePrice: averagePrice,
	})
}

// DeletePortfolioItem removes a portfolio item.
func (r *PortfolioRepository) DeletePortfolioItem(
	ctx context.Context,
	ladderID int64,
	userID int64,
	symbol string,
) error {
	return r.queries.DeletePortfolioItem(ctx, sqlc.DeletePortfolioItemParams{
		LadderID:    ladderID,
		UserID:      userID,
		StockSymbol: symbol,
	})
}
