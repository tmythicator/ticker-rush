package service

import (
	"context"

	pb "github.com/tmythicator/ticker-rush/server/proto/user"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo      UserRepository
	portfolioRepo PortfolioRepository
}

type UserWithPortfolio struct {
	*pb.User
	Portfolio map[string]*pb.PortfolioItem `json:"portfolio"`
}

func NewUserService(userRepo UserRepository, portfolioRepo PortfolioRepository) *UserService {
	return &UserService{
		userRepo:      userRepo,
		portfolioRepo: portfolioRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, email string, password string, firstName string, lastName string) (*pb.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return s.userRepo.CreateUser(ctx, email, string(hashedPassword), firstName, lastName, 10000)
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*pb.User, error) {
	return s.userRepo.GetUser(ctx, id)
}

func (s *UserService) GetUserWithPortfolio(ctx context.Context, id int64) (*UserWithPortfolio, error) {
	user, err := s.userRepo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	portfolio, err := s.portfolioRepo.GetPortfolio(ctx, id)
	if err != nil {
		return nil, err
	}

	userWithPortfolio := &UserWithPortfolio{
		User:      user,
		Portfolio: make(map[string]*pb.PortfolioItem),
	}

	for _, item := range portfolio {
		userWithPortfolio.Portfolio[item.StockSymbol] = &pb.PortfolioItem{
			StockSymbol:  item.StockSymbol,
			Quantity:     item.Quantity,
			AveragePrice: item.AveragePrice,
		}
	}
	return userWithPortfolio, nil
}

func (s *UserService) Authenticate(ctx context.Context, email string, password string) (*pb.User, error) {
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
