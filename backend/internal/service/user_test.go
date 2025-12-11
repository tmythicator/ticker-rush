package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tmythicator/ticker-rush/server/internal/service"
	"github.com/tmythicator/ticker-rush/server/internal/service/mocks"
	pb "github.com/tmythicator/ticker-rush/server/proto/user"

	"golang.org/x/crypto/bcrypt"
)

var (
	ctx = context.Background()
)

func TestUserService_CreateUser(t *testing.T) {
	const password = "test33"
	const email = "abc@gmail.com"
	const startBalance = 10000.0

	mockUserRepo := new(mocks.MockUserRepository)
	mockPortfolioRepo := new(mocks.MockPortfolioRepository)

	expectedUser := &pb.User{
		Id:        1,
		Email:     email,
		FirstName: "John",
		LastName:  "Doe",
	}

	mockUserRepo.On("CreateUser", ctx, email, mock.MatchedBy(func(hashedPassword string) bool {
		return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
	}), expectedUser.FirstName, expectedUser.LastName, startBalance).Return(expectedUser, nil)

	userService := service.NewUserService(mockUserRepo, mockPortfolioRepo)
	user, err := userService.CreateUser(ctx, email, password, expectedUser.FirstName, expectedUser.LastName)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser, user)

	mockUserRepo.AssertExpectations(t)
}

func TestUserService_GetUserWithPortfolio(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockPortfolioRepo := new(mocks.MockPortfolioRepository)
	userService := service.NewUserService(mockUserRepo, mockPortfolioRepo)

	userID := int64(1)
	expectedUser := &pb.User{Id: userID, Email: "test@test.com"}

	expectedPortfolioItems := []*pb.PortfolioItem{
		{StockSymbol: "AAPL", Quantity: 10, AveragePrice: 150.0},
		{StockSymbol: "GOOG", Quantity: 5, AveragePrice: 2000.0},
	}

	mockUserRepo.On("GetUser", ctx, userID).Return(expectedUser, nil)
	mockPortfolioRepo.On("GetPortfolio", ctx, userID).Return(expectedPortfolioItems, nil)

	res, err := userService.GetUserWithPortfolio(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, res.User)
	assert.Len(t, res.Portfolio, len(expectedPortfolioItems))

	assert.Equal(t, expectedPortfolioItems[0], res.Portfolio["AAPL"])
	assert.Equal(t, expectedPortfolioItems[1], res.Portfolio["GOOG"])
}

func TestUserService_Authenticate(t *testing.T) {
	const truePassword = "password123"
	const wrongPassword = "wrongpassword"
	const email = "test@example.com"

	mockUserRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockUserRepo, nil)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(truePassword), bcrypt.DefaultCost)

	expectedUser := &pb.User{
		Id:        1,
		Email:     email,
		FirstName: "John",
		LastName:  "Doe",
	}

	mockUserRepo.On("GetUserByEmail", ctx, email).Return(expectedUser, string(hashedPassword), nil)

	user, err := userService.Authenticate(ctx, email, truePassword)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	_, err = userService.Authenticate(ctx, email, wrongPassword)
	assert.Error(t, err)
}
