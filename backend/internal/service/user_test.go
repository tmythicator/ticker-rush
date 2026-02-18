package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	pb "github.com/tmythicator/ticker-rush/server/internal/proto/user/v1"
	"github.com/tmythicator/ticker-rush/server/internal/service"
	"github.com/tmythicator/ticker-rush/server/internal/service/mocks"
	"golang.org/x/crypto/bcrypt"
)

var ctx = context.Background()

func TestUserService_CreateUser(t *testing.T) {
	const (
		password     = "test134567"
		username     = "newUsername"
		startBalance = 10000.0
	)

	mockUserRepo := new(mocks.MockUserRepository)
	mockPortfolioRepo := new(mocks.MockPortfolioRepository)

	expectedUser := &pb.User{
		Id:        1,
		Username:  username,
		FirstName: "John",
		LastName:  "Doe",
	}

	mockUserRepo.On("CreateUser", ctx, username, mock.MatchedBy(func(hashedPassword string) bool {
		return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
	}), expectedUser.GetFirstName(), expectedUser.GetLastName(), startBalance, "", false).Return(expectedUser, nil)

	userService := service.NewUserService(mockUserRepo, mockPortfolioRepo)
	user, err := userService.CreateUser(
		ctx,
		username,
		password,
		expectedUser.GetFirstName(),
		expectedUser.GetLastName(),
		"",
	)

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
	expectedUser := &pb.User{Id: userID, Username: "userTest1"}

	expectedPortfolioItems := []*pb.PortfolioItem{
		{StockSymbol: "AAPL", Quantity: 10, AveragePrice: 150.0},
		{StockSymbol: "GOOG", Quantity: 5, AveragePrice: 2000.0},
	}

	mockUserRepo.On("GetUser", ctx, userID).Return(expectedUser, nil)
	mockPortfolioRepo.On("GetPortfolio", ctx, userID).Return(expectedPortfolioItems, nil)

	res, err := userService.GetUserWithPortfolio(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.GetId(), res.GetId())
	assert.Equal(t, expectedUser.GetUsername(), res.GetUsername())
	assert.Len(t, res.GetPortfolio(), len(expectedPortfolioItems))

	assert.Equal(t, expectedPortfolioItems[0], res.GetPortfolio()["AAPL"])
	assert.Equal(t, expectedPortfolioItems[1], res.GetPortfolio()["GOOG"])
}

func TestUserService_Authenticate(t *testing.T) {
	const (
		truePassword  = "password123"
		wrongPassword = "wrongpassword"
		username      = "userTest2"
	)

	mockUserRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockUserRepo, nil)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(truePassword), bcrypt.DefaultCost)

	expectedUser := &pb.User{
		Id:        1,
		Username:  username,
		FirstName: "John",
		LastName:  "Doe",
	}

	mockUserRepo.On("GetUserByUsername", ctx, username).Return(expectedUser, string(hashedPassword), nil)

	user, err := userService.Authenticate(ctx, username, truePassword)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	_, err = userService.Authenticate(ctx, username, wrongPassword)
	assert.Error(t, err)
}

func TestUserService_GetPublicProfile(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockPortfolioRepo := new(mocks.MockPortfolioRepository)
	userService := service.NewUserService(mockUserRepo, mockPortfolioRepo)

	t.Run("Public User", func(t *testing.T) {
		username := "publicUser"
		expectedUser := &pb.User{
			Id:        1,
			Username:  username,
			FirstName: "Public",
			LastName:  "User",
			IsPublic:  true,
		}
		expectedPortfolio := []*pb.PortfolioItem{
			{StockSymbol: "AAPL", Quantity: 10, AveragePrice: 150.0},
		}

		mockUserRepo.On("GetUserByUsername", ctx, username).Return(expectedUser, "", nil).Once()
		mockPortfolioRepo.On("GetPortfolio", ctx, expectedUser.Id).Return(expectedPortfolio, nil).Once()

		res, err := userService.GetPublicProfile(ctx, username)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser.Id, res.Id)
		assert.Len(t, res.Portfolio, 1)
		assert.Equal(t, "AAPL", res.Portfolio["AAPL"].StockSymbol)
	})

	t.Run("Private User", func(t *testing.T) {
		username := "privateUser"
		expectedUser := &pb.User{
			Id:       2,
			Username: username,
			IsPublic: false,
		}

		mockUserRepo.On("GetUserByUsername", ctx, username).Return(expectedUser, "", nil).Once()

		_, err := userService.GetPublicProfile(ctx, username)

		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})

	t.Run("User Not Found", func(t *testing.T) {
		username := "unknownUser"

		mockUserRepo.On("GetUserByUsername", ctx, username).Return(nil, "", assert.AnError).Once()

		_, err := userService.GetPublicProfile(ctx, username)

		assert.Error(t, err)
	})
}
