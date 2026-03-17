package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/tmythicator/ticker-rush/backend/internal/gen/sqlc"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/user/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/service"
)

// UserRepository handles user data persistence in PostgreSQL.
type UserRepository struct {
	queries *sqlc.Queries
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		queries: sqlc.New(pool),
	}
}

// GetUser retrieves a user by ID.
func (r *UserRepository) GetUser(ctx context.Context, id int64) (*user.User, error) {
	u, err := r.queries.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return &user.User{
		Id:        u.ID,
		Username:  u.Username,
		CreatedAt: timestamppb.New(u.CreatedAt.Time),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Website:   u.Website,
		IsPublic:  u.IsPublic,
		IsAdmin:   u.IsAdmin,
	}, nil
}

// GetUsers retrieves all users.
func (r *UserRepository) GetUsers(ctx context.Context) ([]*user.User, error) {
	res, err := r.queries.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]*user.User, len(res))

	for i, u := range res {
		users[i] = &user.User{
			Id:        u.ID,
			Username:  u.Username,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			CreatedAt: timestamppb.New(u.CreatedAt.Time),
			Website:   u.Website,
			IsPublic:  u.IsPublic,
			IsAdmin:   u.IsAdmin,
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
func (r *UserRepository) GetUserForUpdate(ctx context.Context, id int64) (*user.User, error) {
	u, err := r.queries.GetUserForUpdate(ctx, id)
	if err != nil {
		return nil, err
	}

	return &user.User{
		Id:        u.ID,
		Username:  u.Username,
		CreatedAt: timestamppb.New(u.CreatedAt.Time),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Website:   u.Website,
		IsPublic:  u.IsPublic,
		IsAdmin:   u.IsAdmin,
	}, nil
}

// UpdateUserProfile updates an existing user's profile.
func (r *UserRepository) UpdateUserProfile(ctx context.Context, u *user.User) error {
	err := r.queries.UpdateUserProfile(ctx, sqlc.UpdateUserProfileParams{
		ID:        u.GetId(),
		FirstName: u.GetFirstName(),
		LastName:  u.GetLastName(),
		Website:   u.GetWebsite(),
		IsPublic:  u.GetIsPublic(),
	})

	return err
}

// UpdateUserBalance updates the user's balance for the given ladder.
func (r *UserRepository) UpdateUserBalance(ctx context.Context, userID int64, ladderID int64, balance float64) error {
	return r.queries.UpdateLadderBalance(ctx, sqlc.UpdateLadderBalanceParams{
		LadderID: ladderID,
		UserID:   userID,
		Balance:  balance,
	})
}

// GetUserBalance retrieves the user's balance for the given ladder.
func (r *UserRepository) GetUserBalance(ctx context.Context, userID int64, ladderID int64) (float64, error) {
	balance, err := r.queries.GetLadderBalance(ctx, sqlc.GetLadderBalanceParams{
		LadderID: ladderID,
		UserID:   userID,
	})
	if err != nil {
		return 0, err
	}

	return balance, nil
}

// CreateUser creates a new user in the database.
func (r *UserRepository) CreateUser(
	ctx context.Context,
	username string,
	hashedPassword string,
	firstName string,
	lastName string,
	website string,
	isPublic bool,
	agbAcceptedAt time.Time,
) (*user.User, error) {
	u, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Username:      username,
		PasswordHash:  hashedPassword,
		FirstName:     firstName,
		LastName:      lastName,
		Website:       website,
		CreatedAt:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
		IsPublic:      isPublic,
		AgbAcceptedAt: pgtype.Timestamptz{Time: agbAcceptedAt, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &user.User{
		Id:        u.ID,
		Username:  u.Username,
		CreatedAt: timestamppb.New(u.CreatedAt.Time),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Website:   u.Website,
		IsPublic:  u.IsPublic,
		IsAdmin:   u.IsAdmin,
	}, nil

}

// GetUserByUsername retrieves a user by username, returning the user and password hash.
func (r *UserRepository) GetUserByUsername(
	ctx context.Context,
	username string,
) (*user.User, string, error) {
	u, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, "", err
	}

	return &user.User{
		Id:        u.ID,
		Username:  u.Username,
		CreatedAt: timestamppb.New(u.CreatedAt.Time),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Website:   u.Website,
		IsPublic:  u.IsPublic,
		IsAdmin:   u.IsAdmin,
	}, u.PasswordHash, nil
}

// GetUserWithPortfolioForActiveLadder retrieves a user and their portfolio for the active ladder.
func (r *UserRepository) GetUserWithPortfolioForActiveLadder(ctx context.Context, userID int64) ([]sqlc.GetUserWithPortfolioForActiveLadderRow, error) {
	return r.queries.GetUserWithPortfolioForActiveLadder(ctx, userID)
}
