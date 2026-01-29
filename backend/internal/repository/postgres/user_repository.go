package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/tmythicator/ticker-rush/server/internal/gen/sqlc"
	pb "github.com/tmythicator/ticker-rush/server/internal/proto/user/v1"
	"github.com/tmythicator/ticker-rush/server/internal/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserRepository handles user data persistence in PostgreSQL.
type UserRepository struct {
	queries *db.Queries
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		queries: db.New(pool),
	}
}

// GetUser retrieves a user by ID.
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

// WithTx returns a new UserRepository that uses the given transaction.
func (r *UserRepository) WithTx(tx service.Transaction) service.UserRepository {
	return &UserRepository{
		queries: r.queries.WithTx(tx.(pgx.Tx)),
	}
}

// GetUserForUpdate retrieves a user by ID with a lock for update.
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

// SaveUser updates an existing user.
func (r *UserRepository) SaveUser(ctx context.Context, user *pb.User) error {
	err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:        user.GetId(),
		Email:     user.GetEmail(),
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		Balance:   user.GetBalance(),
	})

	return err
}

// CreateUser creates a new user in the database.
func (r *UserRepository) CreateUser(
	ctx context.Context,
	email string,
	hashedPassword string,
	firstName string,
	lastName string,
	balance float64,
) (*pb.User, error) {
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

// GetUserByEmail retrieves a user by email, returning the user and password hash.
func (r *UserRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*pb.User, string, error) {
	u, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}

	return &pb.User{
		Id:        u.ID,
		Email:     u.Email,
		Balance:   u.Balance,
		CreatedAt: timestamppb.New(u.CreatedAt.Time),
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}, u.PasswordHash, nil
}
