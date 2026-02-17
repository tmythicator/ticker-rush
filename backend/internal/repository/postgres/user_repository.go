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
		Id:        u.ID,
		Username:  u.Username,
		Balance:   u.Balance,
		CreatedAt: timestamppb.New(u.CreatedAt.Time),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Website:   u.Website,
	}, nil
}

// GetUsers retrieves all users.
func (r *UserRepository) GetUsers(ctx context.Context) ([]*pb.User, error) {
	res, err := r.queries.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]*pb.User, len(res))

	for i, u := range res {
		users[i] = &pb.User{
			Id:        u.ID,
			Username:  u.Username,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Balance:   u.Balance,
			CreatedAt: timestamppb.New(u.CreatedAt.Time),
			Website:   u.Website,
		}
	}

	return users, nil
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
		Username:  u.Username,
		Balance:   u.Balance,
		CreatedAt: timestamppb.New(u.CreatedAt.Time),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Website:   u.Website,
	}, nil
}

// UpdateUserProfile updates an existing user's profile.
func (r *UserRepository) UpdateUserProfile(ctx context.Context, user *pb.User) error {
	err := r.queries.UpdateUserProfile(ctx, db.UpdateUserProfileParams{
		ID:        user.GetId(),
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		Website:   user.GetWebsite(),
	})

	return err
}

// UpdateUserBalance updates the user's balance.
func (r *UserRepository) UpdateUserBalance(ctx context.Context, id int64, balance float64) error {
	return r.queries.UpdateUserBalance(ctx, db.UpdateUserBalanceParams{
		ID:      id,
		Balance: balance,
	})
}

// CreateUser creates a new user in the database.
func (r *UserRepository) CreateUser(
	ctx context.Context,
	username string,
	hashedPassword string,
	firstName string,
	lastName string,
	balance float64,
	website string,
) (*pb.User, error) {
	user, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Username:     username,
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
		Username:  user.Username,
		Balance:   user.Balance,
		CreatedAt: timestamppb.New(user.CreatedAt.Time),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Website:   user.Website,
	}, nil
}

// GetUserByUsername retrieves a user by username, returning the user and password hash.
func (r *UserRepository) GetUserByUsername(
	ctx context.Context,
	username string,
) (*pb.User, string, error) {
	u, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, "", err
	}

	return &pb.User{
		Id:        u.ID,
		Username:  u.Username,
		Balance:   u.Balance,
		CreatedAt: timestamppb.New(u.CreatedAt.Time),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Website:   u.Website,
	}, u.PasswordHash, nil
}
