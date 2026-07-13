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

	userService := service.NewUser(mockUserRepo, mockPortfolioRepo, mockLadderRepo)
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
	userService := service.NewUser(mockUserRepo, mockPortfolioRepo, mockLadderRepo)
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
	userService := service.NewUser(mockUserRepo, mockPortfolioRepo, mockLadderRepo)

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
	userService := service.NewUser(mockUserRepo, mockPortfolioRepo, mockLadderRepo)
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
	userService := service.NewUser(mockUserRepo, nil, mockLadderRepo)

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
	userService := service.NewUser(mockUserRepo, nil, nil)

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
	userService := service.NewUser(mockUserRepo, mockPortfolioRepo, mockLadderRepo)

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

func TestUserService_UpdateUser(t *testing.T) {
	userID := int64(1)

	t.Run("Success with http website", func(t *testing.T) {
		mockUserRepo := new(mocks.MockUserRepository)
		userService := service.NewUser(mockUserRepo, nil, nil)
		existingUser := &domain.User{
			ID:        userID,
			FirstName: "Old",
			LastName:  "Name",
			Website:   "",
			IsPublic:  false,
		}
		mockUserRepo.On("GetUser", ctx, userID).Return(existingUser, nil).Once()
		mockUserRepo.On("UpdateUserProfile", ctx, mock.MatchedBy(func(u *domain.User) bool {
			return u.FirstName == "New" && u.LastName == "Name" && u.Website == "http://example.com" && u.IsPublic == true
		})).Return(nil).Once()

		res, err := userService.UpdateUser(ctx, userID, "New", "Name", "http://example.com", true, nil)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "New", res.FirstName)
		assert.Equal(t, "http://example.com", res.Website)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Success with https website", func(t *testing.T) {
		mockUserRepo := new(mocks.MockUserRepository)
		userService := service.NewUser(mockUserRepo, nil, nil)
		existingUser := &domain.User{
			ID:        userID,
			FirstName: "Old",
			LastName:  "Name",
			Website:   "",
			IsPublic:  false,
		}
		mockUserRepo.On("GetUser", ctx, userID).Return(existingUser, nil).Once()
		mockUserRepo.On("UpdateUserProfile", ctx, mock.MatchedBy(func(u *domain.User) bool {
			return u.Website == "https://example.com"
		})).Return(nil).Once()

		res, err := userService.UpdateUser(ctx, userID, "New", "Name", "https://example.com", true, nil)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Success with empty website", func(t *testing.T) {
		mockUserRepo := new(mocks.MockUserRepository)
		userService := service.NewUser(mockUserRepo, nil, nil)
		existingUser := &domain.User{
			ID:        userID,
			FirstName: "Old",
			LastName:  "Name",
			Website:   "",
			IsPublic:  false,
		}
		mockUserRepo.On("GetUser", ctx, userID).Return(existingUser, nil).Once()
		mockUserRepo.On("UpdateUserProfile", ctx, mock.MatchedBy(func(u *domain.User) bool {
			return u.Website == ""
		})).Return(nil).Once()

		res, err := userService.UpdateUser(ctx, userID, "New", "Name", "", true, nil)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Failure with data URL website", func(t *testing.T) {
		mockUserRepo := new(mocks.MockUserRepository)
		userService := service.NewUser(mockUserRepo, nil, nil)
		existingUser := &domain.User{
			ID:        userID,
			FirstName: "Old",
			LastName:  "Name",
			Website:   "",
			IsPublic:  false,
		}
		mockUserRepo.On("GetUser", ctx, userID).Return(existingUser, nil).Once()

		_, err := userService.UpdateUser(ctx, userID, "New", "Name", "data:text/html,<html></html>", true, nil)
		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrInvalidWebsiteFormat, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Failure with javascript URL website", func(t *testing.T) {
		mockUserRepo := new(mocks.MockUserRepository)
		userService := service.NewUser(mockUserRepo, nil, nil)
		existingUser := &domain.User{
			ID:        userID,
			FirstName: "Old",
			LastName:  "Name",
			Website:   "",
			IsPublic:  false,
		}
		mockUserRepo.On("GetUser", ctx, userID).Return(existingUser, nil).Once()

		_, err := userService.UpdateUser(ctx, userID, "New", "Name", "javascript:alert(1)", true, nil)
		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrInvalidWebsiteFormat, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Failure with no protocol website", func(t *testing.T) {
		mockUserRepo := new(mocks.MockUserRepository)
		userService := service.NewUser(mockUserRepo, nil, nil)
		existingUser := &domain.User{
			ID:        userID,
			FirstName: "Old",
			LastName:  "Name",
			Website:   "",
			IsPublic:  false,
		}
		mockUserRepo.On("GetUser", ctx, userID).Return(existingUser, nil).Once()

		_, err := userService.UpdateUser(ctx, userID, "New", "Name", "google.com", true, nil)
		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrInvalidWebsiteFormat, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Failure with too long website", func(t *testing.T) {
		mockUserRepo := new(mocks.MockUserRepository)
		userService := service.NewUser(mockUserRepo, nil, nil)
		existingUser := &domain.User{
			ID:        userID,
			FirstName: "Old",
			LastName:  "Name",
			Website:   "",
			IsPublic:  false,
		}
		mockUserRepo.On("GetUser", ctx, userID).Return(existingUser, nil).Once()

		longWebsite := "https://example.com/" + string(make([]byte, 201))
		_, err := userService.UpdateUser(ctx, userID, "New", "Name", longWebsite, true, nil)
		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrInvalidWebsiteFormat, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Failure with name required", func(t *testing.T) {
		mockUserRepo := new(mocks.MockUserRepository)
		userService := service.NewUser(mockUserRepo, nil, nil)
		existingUser := &domain.User{
			ID:        userID,
			FirstName: "Old",
			LastName:  "Name",
			Website:   "",
			IsPublic:  false,
		}
		mockUserRepo.On("GetUser", ctx, userID).Return(existingUser, nil).Once()

		_, err := userService.UpdateUser(ctx, userID, "", "Name", "https://example.com", true, nil)
		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrNameRequired, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Failure with profanity name", func(t *testing.T) {
		mockUserRepo := new(mocks.MockUserRepository)
		userService := service.NewUser(mockUserRepo, nil, nil)
		existingUser := &domain.User{
			ID:        userID,
			FirstName: "Old",
			LastName:  "Name",
			Website:   "",
			IsPublic:  false,
		}
		mockUserRepo.On("GetUser", ctx, userID).Return(existingUser, nil).Once()

		_, err := userService.UpdateUser(ctx, userID, "fuck", "Name", "https://example.com", true, nil)
		assert.Error(t, err)
		assert.Equal(t, apperrors.ErrProfanityDetected, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Success partial update with paths", func(t *testing.T) {
		mockUserRepo := new(mocks.MockUserRepository)
		userService := service.NewUser(mockUserRepo, nil, nil)
		existingUser := &domain.User{
			ID:        userID,
			FirstName: "Old",
			LastName:  "Name",
			Website:   "http://old.com",
			IsPublic:  false,
		}
		mockUserRepo.On("GetUser", ctx, userID).Return(existingUser, nil).Once()
		mockUserRepo.On("UpdateUserProfile", ctx, mock.MatchedBy(func(u *domain.User) bool {
			return u.FirstName == "Old" && u.LastName == "Name" && u.Website == "http://new.com" && u.IsPublic == false
		})).Return(nil).Once()

		res, err := userService.UpdateUser(ctx, userID, "ShouldNotUpdate", "ShouldNotUpdate", "http://new.com", true, []string{"website"})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "Old", res.FirstName)
		assert.Equal(t, "http://new.com", res.Website)
		assert.False(t, res.IsPublic)
		mockUserRepo.AssertExpectations(t)
	})
}
