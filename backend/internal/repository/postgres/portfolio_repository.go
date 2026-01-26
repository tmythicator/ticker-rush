package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/tmythicator/ticker-rush/server/internal/gen/sqlc"
	pb "github.com/tmythicator/ticker-rush/server/internal/proto/user"
	"github.com/tmythicator/ticker-rush/server/internal/service"
)

type PortfolioRepository struct {
	queries *db.Queries
}

func NewPortfolioRepository(pool *pgxpool.Pool) *PortfolioRepository {
	return &PortfolioRepository{
		queries: db.New(pool),
	}
}

func (r *PortfolioRepository) WithTx(tx service.Transaction) service.PortfolioRepository {
	return &PortfolioRepository{
		queries: r.queries.WithTx(tx.(pgx.Tx)),
	}
}

func (r *PortfolioRepository) GetPortfolio(
	ctx context.Context,
	userID int64,
) ([]*pb.PortfolioItem, error) {
	items, err := r.queries.GetPortfolio(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*pb.PortfolioItem, len(items))
	for i, item := range items {
		result[i] = &pb.PortfolioItem{
			StockSymbol:  item.StockSymbol,
			Quantity:     item.Quantity,
			AveragePrice: item.AveragePrice,
		}
	}

	return result, nil
}

func (r *PortfolioRepository) GetPortfolioItem(
	ctx context.Context,
	userID int64,
	symbol string,
) (*pb.PortfolioItem, error) {
	item, err := r.queries.GetPortfolioItem(ctx, db.GetPortfolioItemParams{
		UserID:      userID,
		StockSymbol: symbol,
	})
	if err != nil {
		return nil, err
	}

	return &pb.PortfolioItem{
		StockSymbol:  item.StockSymbol,
		Quantity:     item.Quantity,
		AveragePrice: item.AveragePrice,
	}, nil
}

func (r *PortfolioRepository) GetPortfolioItemForUpdate(
	ctx context.Context,
	userID int64,
	symbol string,
) (*pb.PortfolioItem, error) {
	item, err := r.queries.GetPortfolioItemForUpdate(ctx, db.GetPortfolioItemForUpdateParams{
		UserID:      userID,
		StockSymbol: symbol,
	})
	if err != nil {
		return nil, err
	}

	return &pb.PortfolioItem{
		StockSymbol:  item.StockSymbol,
		Quantity:     item.Quantity,
		AveragePrice: item.AveragePrice,
	}, nil
}

func (r *PortfolioRepository) SetPortfolioItem(
	ctx context.Context,
	userID int64,
	symbol string,
	quantity float64,
	averagePrice float64,
) error {
	return r.queries.SetPortfolioItem(ctx, db.SetPortfolioItemParams{
		UserID:       userID,
		StockSymbol:  symbol,
		Quantity:     quantity,
		AveragePrice: averagePrice,
	})
}

func (r *PortfolioRepository) DeletePortfolioItem(
	ctx context.Context,
	userID int64,
	symbol string,
) error {
	return r.queries.DeletePortfolioItem(ctx, db.DeletePortfolioItemParams{
		UserID:      userID,
		StockSymbol: symbol,
	})
}
