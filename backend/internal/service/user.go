package service

import (
	"context"
	"errors"
	"slices"
	"time"

	"regexp"

	goaway "github.com/TwiN/go-away"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/tmythicator/ticker-rush/backend/internal/apperrors"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/portfolio/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/user/v1"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
)

// UserService handles user-related business logic.
type UserService struct {
	userRepo      UserRepository
	portfolioRepo PortfolioRepository
	ladderRepo    LadderRepository
}

// NewUserService creates a new instance of UserService.
func NewUserService(userRepo UserRepository, portfolioRepo PortfolioRepository, ladderRepo LadderRepository) *UserService {
	return &UserService{
		userRepo:      userRepo,
		portfolioRepo: portfolioRepo,
		ladderRepo:    ladderRepo,
	}
}

// CreateUser registers a new user.
func (s *UserService) CreateUser(
	ctx context.Context,
	username string,
	password string,
	firstName string,
	lastName string,
	agbAccepted bool,
) (*user.User, error) {
	// Require AGB Acceptance
	if !agbAccepted {
		return nil, apperrors.ErrAGBNotAccepted
	}

	// 1. Validate Username Format
	if !usernameRegex.MatchString(username) {
		return nil, apperrors.ErrInvalidUsernameFormat
	}

	// Validate Password Length
	if len(password) < 8 {
		return nil, apperrors.ErrPasswordTooShort
	}

	// Validate Names
	if len(firstName) == 0 || len(lastName) == 0 {
		return nil, apperrors.ErrNameRequired
	}

	// 2. Profanity Check
	if goaway.IsProfane(username) || goaway.IsProfane(firstName) || goaway.IsProfane(lastName) {
		return nil, apperrors.ErrProfanityDetected
	}

	// 3. Reserved usernames
	blockedNames := []string{"admin", "administrator", "system", "mod", "moderator", "support", "help"}
	if slices.Contains(blockedNames, username) {
		return nil, apperrors.ErrUsernameNotAllowed
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return s.userRepo.CreateUser(ctx, username, string(hashedPassword), firstName, lastName, "", false, time.Now())
}

// GetUser retrieves a user by ID.
func (s *UserService) GetUser(ctx context.Context, id int64) (*user.User, error) {
	return s.userRepo.GetUser(ctx, id)
}

// GetUserWithPortfolio retrieves a user and their portfolio for the active ladder.
func (s *UserService) GetUserWithPortfolio(ctx context.Context, id int64) (*user.User, error) {
	rows, err := s.userRepo.GetUserWithPortfolioForActiveLadder(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, apperrors.ErrUserNotFound
	}
	firstRow := rows[0]

	fetchedUser := &user.User{
		Id:              firstRow.UserID,
		Username:        firstRow.Username,
		FirstName:       firstRow.FirstName,
		LastName:        firstRow.LastName,
		Website:         firstRow.Website,
		IsPublic:        firstRow.IsPublic,
		IsAdmin:         firstRow.IsAdmin,
		IsBanned:        firstRow.IsBanned,
		CreatedAt:       timestamppb.New(firstRow.CreatedAt.Time),
		Balance:         firstRow.Balance,
		Portfolio:       make(map[string]*portfolio.PortfolioItem),
		IsParticipating: firstRow.IsParticipating,
	}

	for _, row := range rows {
		if row.StockSymbol.Valid && row.LadderID > 0 {
			fetchedUser.Portfolio[row.StockSymbol.String] = &portfolio.PortfolioItem{
				StockSymbol:  row.StockSymbol.String,
				Quantity:     row.Quantity,
				AveragePrice: row.AveragePrice,
			}
		}
	}

	return fetchedUser, nil
}

// Authenticate checks user credentials.
func (s *UserService) Authenticate(
	ctx context.Context,
	username string,
	password string,
) (*user.User, error) {
	user, passwordHash, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates a user's profile.
func (s *UserService) UpdateUser(
	ctx context.Context,
	id int64,
	firstName string,
	lastName string,
	website string,
	isPublic bool,
) (*user.User, error) {
	// Validate Names
	if len(firstName) == 0 || len(lastName) == 0 {
		return nil, apperrors.ErrNameRequired
	}

	// Profanity Check
	if goaway.IsProfane(firstName) || goaway.IsProfane(lastName) {
		return nil, apperrors.ErrProfanityDetected
	}

	// Get existing user to preserve other fields
	existingUser, err := s.userRepo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	existingUser.FirstName = firstName
	existingUser.LastName = lastName
	existingUser.Website = website
	existingUser.IsPublic = isPublic

	if err := s.userRepo.UpdateUserProfile(ctx, existingUser); err != nil {
		return nil, err
	}

	return existingUser, nil
}

// GetPublicProfile retrieves a user's public profile if enabled, for the active ladder.
func (s *UserService) GetPublicProfile(ctx context.Context, username string) (*user.User, error) {
	targetUser, _, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}

		return nil, err
	}

	if !targetUser.IsPublic {
		return nil, apperrors.ErrUserNotFound
	}

	return s.GetUserWithPortfolio(ctx, targetUser.Id)
}
