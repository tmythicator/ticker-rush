package service

import (
	"context"

	"github.com/tmythicator/ticker-rush/server/internal/proto/user"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo      UserRepository
	portfolioRepo PortfolioRepository
}

func NewUserService(userRepo UserRepository, portfolioRepo PortfolioRepository) *UserService {
	return &UserService{
		userRepo:      userRepo,
		portfolioRepo: portfolioRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, email string, password string, firstName string, lastName string) (*user.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return s.userRepo.CreateUser(ctx, email, string(hashedPassword), firstName, lastName, 10000)
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*user.User, error) {
	return s.userRepo.GetUser(ctx, id)
}

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
		fetchedUser.Portfolio[item.StockSymbol] = &user.PortfolioItem{
			StockSymbol:  item.StockSymbol,
			Quantity:     item.Quantity,
			AveragePrice: item.AveragePrice,
		}
	}
	return fetchedUser, nil
}

func (s *UserService) Authenticate(ctx context.Context, email string, password string) (*user.User, error) {
	user, passwordHash, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}
