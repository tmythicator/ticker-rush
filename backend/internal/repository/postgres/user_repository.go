package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/tmythicator/ticker-rush/backend/internal/domain"
	"github.com/tmythicator/ticker-rush/backend/internal/gen/sqlc"
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
func (r *UserRepository) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	u, err := r.queries.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        u.ID,
		Username:  u.Username,
		CreatedAt: u.CreatedAt.Time,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Website:   u.Website,
		IsPublic:  u.IsPublic,
		IsAdmin:   u.IsAdmin,
	}, nil
}

// GetUsers retrieves all users.
func (r *UserRepository) GetUsers(ctx context.Context) ([]*domain.User, error) {
	res, err := r.queries.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]*domain.User, len(res))

	for i, u := range res {
		users[i] = &domain.User{
			ID:        u.ID,
			Username:  u.Username,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			CreatedAt: u.CreatedAt.Time,
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
func (r *UserRepository) GetUserForUpdate(ctx context.Context, id int64) (*domain.User, error) {
	u, err := r.queries.GetUserForUpdate(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        u.ID,
		Username:  u.Username,
		CreatedAt: u.CreatedAt.Time,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Website:   u.Website,
		IsPublic:  u.IsPublic,
		IsAdmin:   u.IsAdmin,
	}, nil
}

// UpdateUserProfile updates an existing user's profile.
func (r *UserRepository) UpdateUserProfile(ctx context.Context, u *domain.User) error {
	err := r.queries.UpdateUserProfile(ctx, sqlc.UpdateUserProfileParams{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Website:   u.Website,
		IsPublic:  u.IsPublic,
	})

	return err
}

// UpdateUserBalance updates the user's balance for the given ladder.
func (r *UserRepository) UpdateUserBalance(ctx context.Context, userID int64, ladderID int64, balance decimal.Decimal) error {
	return r.queries.UpdateLadderParticipantBalance(ctx, sqlc.UpdateLadderParticipantBalanceParams{
		LadderID: ladderID,
		UserID:   userID,
		Balance:  balance,
	})
}

// GetUserBalance retrieves the user's balance for the given ladder.
func (r *UserRepository) GetUserBalance(ctx context.Context, userID int64, ladderID int64) (decimal.Decimal, error) {
	balance, err := r.queries.GetLadderParticipantBalance(ctx, sqlc.GetLadderParticipantBalanceParams{
		LadderID: ladderID,
		UserID:   userID,
	})
	if err != nil {
		return decimal.Zero, err
	}

	return balance, nil
}

// CreateUser creates a new user in the database.
func (r *UserRepository) CreateUser(
	ctx context.Context,
	params service.CreateUserParams,
) (*domain.User, error) {
	u, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Username:      params.Username,
		PasswordHash:  params.PasswordHash,
		FirstName:     params.FirstName,
		LastName:      params.LastName,
		Website:       params.Website,
		CreatedAt:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
		IsPublic:      params.IsPublic,
		AgbAcceptedAt: pgtype.Timestamptz{Time: params.AgbAcceptedAt, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        u.ID,
		Username:  u.Username,
		CreatedAt: u.CreatedAt.Time,
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
) (*domain.User, string, error) {
	u, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, "", err
	}

	return &domain.User{
		ID:        u.ID,
		Username:  u.Username,
		CreatedAt: u.CreatedAt.Time,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Website:   u.Website,
		IsPublic:  u.IsPublic,
		IsAdmin:   u.IsAdmin,
	}, u.PasswordHash, nil
}

// GetUserWithPortfolioForActiveLadder retrieves a user and their portfolio for the active ladder.
func (r *UserRepository) GetUserWithPortfolioForActiveLadder(ctx context.Context, userID int64) (*domain.User, error) {
	rows, err := r.queries.GetUserWithPortfolioForActiveLadder(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, pgx.ErrNoRows
	}
	firstRow := rows[0]

	fetchedUser := &domain.User{
		ID:              firstRow.UserID,
		Username:        firstRow.Username,
		FirstName:       firstRow.FirstName,
		LastName:        firstRow.LastName,
		Website:         firstRow.Website,
		IsPublic:        firstRow.IsPublic,
		IsAdmin:         firstRow.IsAdmin,
		IsBanned:        firstRow.IsBanned,
		CreatedAt:       firstRow.CreatedAt.Time,
		Balance:         firstRow.Balance,
		Portfolio:       make(map[string]domain.PortfolioItem),
		IsParticipating: firstRow.IsParticipating,
	}

	for _, row := range rows {
		if row.StockSymbol.Valid && row.LadderID > 0 {
			fetchedUser.Portfolio[row.StockSymbol.String] = domain.PortfolioItem{
				StockSymbol:  row.StockSymbol.String,
				Quantity:     decimal.NewFromFloat(row.Quantity),
				AveragePrice: decimal.NewFromFloat(row.AveragePrice),
			}
		}
	}

	return fetchedUser, nil
}
