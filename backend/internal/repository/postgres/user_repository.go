package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/tmythicator/ticker-rush/server/internal/gen/sqlc"
	pb "github.com/tmythicator/ticker-rush/server/proto/user"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserRepository struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		queries: db.New(pool),
		pool:    pool,
	}
}

func (r *UserRepository) GetUser(ctx context.Context, id int64) (*pb.User, error) {
	u, err := r.queries.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	items, err := r.queries.GetPortfolio(ctx, id)
	if err != nil {
		return nil, err
	}

	portfolio := make(map[string]int32)
	for _, item := range items {
		portfolio[item.StockSymbol] = item.Quantity
	}

	return &pb.User{
		Id:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Balance:      u.Balance,
		Portfolio:    portfolio,
		CreatedAt:    timestamppb.New(u.CreatedAt.Time),
	}, nil
}

func (r *UserRepository) SaveUser(ctx context.Context, user *pb.User) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	err = qtx.UpsertUser(ctx, db.UpsertUserParams{
		ID:           user.Id,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Balance:      user.Balance,
		CreatedAt:    pgtype.Timestamptz{Time: user.CreatedAt.AsTime(), Valid: true},
	})
	if err != nil {
		return err
	}

	for symbol, quantity := range user.Portfolio {
		err := qtx.SetPortfolioItem(ctx, db.SetPortfolioItemParams{
			UserID:       user.Id,
			StockSymbol:  symbol,
			Quantity:     quantity,
			AveragePrice: 0,
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *UserRepository) CreateUser(ctx context.Context, id int64, password, email string) (*pb.User, error) {
	exists, err := r.queries.CheckUserExists(ctx, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &pb.User{
		Id:           id,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Balance:      10000.0, // TODO: Make configurable each ladder run
		Portfolio:    make(map[string]int32),
		CreatedAt:    timestamppb.Now(),
	}

	_, err = r.queries.CreateUser(ctx, db.CreateUserParams{
		ID:           user.Id,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Balance:      user.Balance,
		CreatedAt:    pgtype.Timestamptz{Time: user.CreatedAt.AsTime(), Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
