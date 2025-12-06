package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
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

	return &pb.User{
		Id:    u.ID,
		Email: u.Email,
		// PasswordHash is excluded from the query for security
		Balance:   u.Balance,
		CreatedAt: timestamppb.New(u.CreatedAt.Time),
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}, nil
}

func (r *UserRepository) WithTx(tx pgx.Tx) *UserRepository {
	return &UserRepository{
		queries: r.queries.WithTx(tx),
		pool:    r.pool,
	}
}

func (r *UserRepository) GetUserForUpdate(ctx context.Context, id int64) (*pb.User, error) {
	u, err := r.queries.GetUserForUpdate(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id:        u.ID,
		Email:     u.Email,
		Balance:   u.Balance,
		CreatedAt: timestamppb.New(u.CreatedAt.Time),
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}, nil
}

func (r *UserRepository) SaveUser(ctx context.Context, user *pb.User) error {
	// Note: When using WithTx, r.queries is already bound to the transaction.
	// If r.pool is used here to begin a transaction, it would be a nested transaction
	// or independent if not using the tx.
	// Since we want this to be part of the external transaction when checking from main,
	// we should rely on r.queries being set correctly via WithTx.
	// However, if called without WithTx, r.queries uses the db (pool).

	// But wait, UpdateUser is just a query. The previous implementation of SaveUser
	// handled the PORTFOLIO deletions too. Now SaveUser ONLY updates the user balance/details.
	// So we don't need a transaction inside SaveUser anymore if we assume the caller
	// orchestrates the transaction for atomic User+Portfolio updates.

	err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:        user.Id,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Balance:   user.Balance,
	})
	return err
}

func (r *UserRepository) CreateUser(ctx context.Context, email string, password string, firstName string, lastName string) (*pb.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: string(hashedPassword),
		FirstName:    firstName,
		LastName:     lastName,
		Balance:      10000.0,
		CreatedAt:    pgtype.Timestamptz{Time: timestamppb.Now().AsTime(), Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		Balance:   user.Balance,
		CreatedAt: timestamppb.New(user.CreatedAt.Time),
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}
