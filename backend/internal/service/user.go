package service

import (
	"context"
	"slices"

	"regexp"

	goaway "github.com/TwiN/go-away"
	"github.com/tmythicator/ticker-rush/server/internal/apperrors"
	"github.com/tmythicator/ticker-rush/server/internal/proto/user/v1"
	"golang.org/x/crypto/bcrypt"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
)

// UserService handles user-related business logic.
type UserService struct {
	userRepo      UserRepository
	portfolioRepo PortfolioRepository
}

// NewUserService creates a new instance of UserService.
func NewUserService(userRepo UserRepository, portfolioRepo PortfolioRepository) *UserService {
	return &UserService{
		userRepo:      userRepo,
		portfolioRepo: portfolioRepo,
	}
}

// CreateUser registers a new user.
func (s *UserService) CreateUser(
	ctx context.Context,
	username string,
	password string,
	firstName string,
	lastName string,
	website string,
) (*user.User, error) {
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

	return s.userRepo.CreateUser(ctx, username, string(hashedPassword), firstName, lastName, 10000, website)
}

// GetUser retrieves a user by ID.
func (s *UserService) GetUser(ctx context.Context, id int64) (*user.User, error) {
	return s.userRepo.GetUser(ctx, id)
}

// GetUserWithPortfolio retrieves a user and their portfolio.
func (s *UserService) GetUserWithPortfolio(ctx context.Context, id int64) (*user.User, error) {
	fetchedUser, err := s.userRepo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	portfolio, err := s.portfolioRepo.GetPortfolio(ctx, id)
	if err != nil {
		return nil, err
	}

	if fetchedUser.Portfolio == nil {
		fetchedUser.Portfolio = make(map[string]*user.PortfolioItem)
	}

	for _, item := range portfolio {
		fetchedUser.Portfolio[item.GetStockSymbol()] = &user.PortfolioItem{
			StockSymbol:  item.GetStockSymbol(),
			Quantity:     item.GetQuantity(),
			AveragePrice: item.GetAveragePrice(),
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
