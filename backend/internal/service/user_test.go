package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/tmythicator/ticker-rush/backend/internal/apperrors"
	"github.com/tmythicator/ticker-rush/backend/internal/domain"
	"github.com/tmythicator/ticker-rush/backend/internal/service"
	"github.com/tmythicator/ticker-rush/backend/internal/service/mocks"
)

var ctx = context.Background()

func TestUserService_CreateUser(t *testing.T) {
	const (
		password = "test134567"
		username = "newUsername"
	)

	mockUserRepo := new(mocks.MockUserRepository)
	mockPortfolioRepo := new(mocks.MockPortfolioRepository)
	mockLadderRepo := new(mocks.MockLadderRepository)

	expectedUser := &domain.User{
		ID:        1,
		Username:  username,
		FirstName: "John",
		LastName:  "Doe",
	}

	mockUserRepo.On("CreateUser", ctx, mock.MatchedBy(func(params service.CreateUserParams) bool {
		return params.Username == username &&
			bcrypt.CompareHashAndPassword([]byte(params.PasswordHash), []byte(password)) == nil &&
			params.FirstName == expectedUser.FirstName &&
			params.LastName == expectedUser.LastName &&
			params.Website == "" &&
			!params.IsPublic
	})).Return(expectedUser, nil)

	userService := service.NewUserService(mockUserRepo, mockPortfolioRepo, mockLadderRepo)
	user, err := userService.CreateUser(
		ctx,
		username,
		password,
		expectedUser.FirstName,
		expectedUser.LastName,
		true,
	)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser, user)

	mockUserRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_AgbNotAccepted(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockPortfolioRepo := new(mocks.MockPortfolioRepository)
	mockLadderRepo := new(mocks.MockLadderRepository)
	userService := service.NewUserService(mockUserRepo, mockPortfolioRepo, mockLadderRepo)
	user, err := userService.CreateUser(
		ctx,
		"username",
		"password",
		"First",
		"Last",
		false,
	)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, apperrors.ErrAGBNotAccepted, err)

	mockUserRepo.AssertNotCalled(t, "CreateUser")
}

func TestUserService_CreateUser_PasswordTooLong(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockPortfolioRepo := new(mocks.MockPortfolioRepository)
	mockLadderRepo := new(mocks.MockLadderRepository)
	userService := service.NewUserService(mockUserRepo, mockPortfolioRepo, mockLadderRepo)

	// 73 chars
	longPassword := "0123456789012345678901234567890123456789012345678901234567890123456789012"
	user, err := userService.CreateUser(
		ctx,
		"username",
		longPassword,
		"First",
		"Last",
		true,
	)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, apperrors.ErrPasswordTooLong, err)

	mockUserRepo.AssertNotCalled(t, "CreateUser")
}

func TestUserService_GetUserWithPortfolio(t *testing.T) {
	usernameTest := "userTest1"
	mockUserRepo := new(mocks.MockUserRepository)
	mockPortfolioRepo := new(mocks.MockPortfolioRepository)
	mockLadderRepo := new(mocks.MockLadderRepository)
	userService := service.NewUserService(mockUserRepo, mockPortfolioRepo, mockLadderRepo)
	userID := int64(1)
	now := time.Now()

	expectedUser := &domain.User{
		ID:        userID,
		Username:  usernameTest,
		FirstName: "Test",
		LastName:  "User",
		CreatedAt: now,
		Balance:   decimal.NewFromFloat(100.0),
		Portfolio: map[string]domain.PortfolioItem{
			"AAPL": {
				StockSymbol:  "AAPL",
				Quantity:     decimal.NewFromFloat(10.0),
				AveragePrice: decimal.NewFromFloat(150.0),
			},
			"GOOG": {
				StockSymbol:  "GOOG",
				Quantity:     decimal.NewFromFloat(5.0),
				AveragePrice: decimal.NewFromFloat(2000.0),
			},
		},
	}

	mockUserRepo.On("GetUserWithPortfolioForActiveLadder", ctx, userID).Return(expectedUser, nil)

	res, err := userService.GetUserWithPortfolio(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, userID, res.ID)
	assert.Equal(t, usernameTest, res.Username)
	assert.Len(t, res.Portfolio, 2)

	assert.Equal(t, float64(10), res.Portfolio["AAPL"].Quantity.InexactFloat64())
	assert.Equal(t, float64(5), res.Portfolio["GOOG"].Quantity.InexactFloat64())
}

func TestUserService_Authenticate(t *testing.T) {
	const (
		truePassword  = "password123"
		wrongPassword = "wrongpassword"
		username      = "userTest2"
	)

	mockUserRepo := new(mocks.MockUserRepository)
	mockLadderRepo := new(mocks.MockLadderRepository)
	userService := service.NewUserService(mockUserRepo, nil, mockLadderRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(truePassword), bcrypt.DefaultCost)

	expectedUser := &domain.User{
		ID:        1,
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

func TestUserService_Authenticate_PasswordTooLong(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	userService := service.NewUserService(mockUserRepo, nil, nil)

	// 73 chars
	longPassword := "0123456789012345678901234567890123456789012345678901234567890123456789012"

	_, err := userService.Authenticate(ctx, "username", longPassword)
	assert.Error(t, err)
	assert.Equal(t, apperrors.ErrPasswordTooLong, err)

	mockUserRepo.AssertNotCalled(t, "GetUserByUsername")
}

func TestUserService_GetPublicProfile(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockPortfolioRepo := new(mocks.MockPortfolioRepository)
	mockLadderRepo := new(mocks.MockLadderRepository)
	userService := service.NewUserService(mockUserRepo, mockPortfolioRepo, mockLadderRepo)

	t.Run("Public User", func(t *testing.T) {
		username := "publicUser"
		expectedUser := &domain.User{
			ID:        1,
			Username:  username,
			FirstName: "Public",
			LastName:  "User",
			IsPublic:  true,
		}
		mockUserRepo.On("GetUserByUsername", ctx, username).Return(expectedUser, "", nil).Once()

		expectedUserWithPortfolio := &domain.User{
			ID:              expectedUser.ID,
			Username:        username,
			FirstName:       "Public",
			LastName:        "User",
			CreatedAt:       time.Now(),
			IsPublic:        true,
			Balance:         decimal.NewFromFloat(200.0),
			IsParticipating: true,
			Portfolio: map[string]domain.PortfolioItem{
				"AAPL": {
					StockSymbol:  "AAPL",
					Quantity:     decimal.NewFromFloat(10.0),
					AveragePrice: decimal.NewFromFloat(150.0),
				},
			},
		}

		mockUserRepo.On("GetUserWithPortfolioForActiveLadder", ctx, expectedUser.ID).Return(expectedUserWithPortfolio, nil).Once()

		res, err := userService.GetPublicProfile(ctx, username)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, res.ID)
		assert.Len(t, res.Portfolio, 1)
		assert.Equal(t, "AAPL", res.Portfolio["AAPL"].StockSymbol)
	})

	t.Run("Private User", func(t *testing.T) {
		username := "privateUser"
		expectedUser := &domain.User{
			ID:       2,
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
