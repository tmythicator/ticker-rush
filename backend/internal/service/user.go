package service

import (
	"context"
	"errors"
	"net/url"
	"regexp"
	"slices"
	"time"

	goaway "github.com/TwiN/go-away"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/tmythicator/ticker-rush/backend/internal/apperrors"
	"github.com/tmythicator/ticker-rush/backend/internal/domain"
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
) (*domain.User, error) {
	if err := validateUserParams(username, password, firstName, lastName, agbAccepted); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return s.userRepo.CreateUser(ctx, CreateUserParams{
		Username:      username,
		PasswordHash:  string(hashedPassword),
		FirstName:     firstName,
		LastName:      lastName,
		Website:       "",
		IsPublic:      false,
		AgbAcceptedAt: time.Now(),
	})
}

// GetUser retrieves a user by ID.
func (s *UserService) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	return s.userRepo.GetUser(ctx, id)
}

// GetUserWithPortfolio retrieves a user and their portfolio for the active ladder.
func (s *UserService) GetUserWithPortfolio(ctx context.Context, id int64) (*domain.User, error) {
	return s.userRepo.GetUserWithPortfolioForActiveLadder(ctx, id)
}

// Authenticate checks user credentials.
func (s *UserService) Authenticate(
	ctx context.Context,
	username string,
	password string,
) (*domain.User, error) {
	if len(password) > 72 {
		return nil, apperrors.ErrPasswordTooLong
	}
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
) (*domain.User, error) {
	// Validate Names
	if len(firstName) == 0 || len(lastName) == 0 {
		return nil, apperrors.ErrNameRequired
	}

	// Profanity Check
	if goaway.IsProfane(firstName) || goaway.IsProfane(lastName) {
		return nil, apperrors.ErrProfanityDetected
	}

	// Website Check
	if website != "" {
		if len(website) > 200 {
			return nil, apperrors.ErrInvalidWebsiteFormat
		}
		u, err := url.ParseRequestURI(website)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
			return nil, apperrors.ErrInvalidWebsiteFormat
		}
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
func (s *UserService) GetPublicProfile(ctx context.Context, username string) (*domain.User, error) {
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

	return s.GetUserWithPortfolio(ctx, targetUser.ID)
}

func validateUserParams(username, password, firstName, lastName string, agbAccepted bool) error {
	if !agbAccepted {
		return apperrors.ErrAGBNotAccepted
	}

	if !usernameRegex.MatchString(username) {
		return apperrors.ErrInvalidUsernameFormat
	}

	if len(password) < 8 {
		return apperrors.ErrPasswordTooShort
	}
	if len(password) > 72 {
		return apperrors.ErrPasswordTooLong
	}

	if len(firstName) == 0 || len(lastName) == 0 {
		return apperrors.ErrNameRequired
	}

	if goaway.IsProfane(username) || goaway.IsProfane(firstName) || goaway.IsProfane(lastName) {
		return apperrors.ErrProfanityDetected
	}

	blockedNames := []string{"admin", "administrator", "system", "mod", "moderator", "support", "help"}
	if slices.Contains(blockedNames, username) {
		return apperrors.ErrUsernameNotAllowed
	}

	return nil
}
