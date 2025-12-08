package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/tmythicator/ticker-rush/server/internal/gen/sqlc"
	pb "github.com/tmythicator/ticker-rush/server/proto/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserRepository struct {
	queries *db.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		queries: db.New(pool),
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
	err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:        user.Id,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Balance:   user.Balance,
	})
	return err
}

func (r *UserRepository) CreateUser(ctx context.Context, email string, hashedPassword string, firstName string, lastName string, balance float64) (*pb.User, error) {
	user, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: hashedPassword,
		FirstName:    firstName,
		LastName:     lastName,
		Balance:      balance,
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
