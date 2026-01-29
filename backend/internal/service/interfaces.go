package service

import (
	"context"

	pb "github.com/tmythicator/ticker-rush/server/internal/proto/user/v1"
)

// Transaction represents a database transaction.
type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// Transactor is responsible for creating transactions.
type Transactor interface {
	Begin(ctx context.Context) (Transaction, error)
}

// UserRepository defines the interface for user persistence.
type UserRepository interface {
	GetUser(ctx context.Context, id int64) (*pb.User, error)
	GetUserByEmail(
		ctx context.Context,
		email string,
	) (*pb.User, string, error) // Returns user, hash, error
	CreateUser(
		ctx context.Context,
		email string,
		hashedPassword string,
		firstName string,
		lastName string,
		balance float64,
	) (*pb.User, error)

	GetUserForUpdate(ctx context.Context, id int64) (*pb.User, error)
	SaveUser(ctx context.Context, user *pb.User) error
	WithTx(tx Transaction) UserRepository
}

// PortfolioRepository defines the interface for portfolio persistence.
type PortfolioRepository interface {
	GetPortfolio(ctx context.Context, userID int64) ([]*pb.PortfolioItem, error)
	GetPortfolioItem(ctx context.Context, userID int64, symbol string) (*pb.PortfolioItem, error)

	GetPortfolioItemForUpdate(
		ctx context.Context,
		userID int64,
		symbol string,
	) (*pb.PortfolioItem, error)
	SetPortfolioItem(
		ctx context.Context,
		userID int64,
		symbol string,
		quantity float64,
		averagePrice float64,
	) error
	DeletePortfolioItem(ctx context.Context, userID int64, symbol string) error
	WithTx(tx Transaction) PortfolioRepository
}
