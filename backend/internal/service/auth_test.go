package service_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	pb "github.com/tmythicator/ticker-rush/server/internal/proto/user/v1"
	"github.com/tmythicator/ticker-rush/server/internal/service"
)

func TestAuthService_GenerateToken(t *testing.T) {
	const (
		testUsername = "testuser"
		testUserID   = int64(123)
		secret       = "test-secret"
	)

	user := &pb.User{
		Id:       testUserID,
		Username: testUsername,
	}

	tokenString, err := service.GenerateToken(user, secret)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, float64(testUserID), claims["user_id"])
	assert.Equal(t, testUsername, claims["username"])
}

func TestAuthService_ValidateToken(t *testing.T) {
	const secret = "test-secret"

	t.Run("Valid Token", func(t *testing.T) {
		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &service.Claims{
			UserID:   456,
			Username: "testuser",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Issuer:    "ticker-rush",
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			t.Fatalf("Failed to sign token: %v", err)
		}

		// Validate the token
		parsedClaims, err := service.ValidateToken(tokenString, secret)
		if err != nil {
			t.Fatalf("ValidateToken failed: %v", err)
		}

		if parsedClaims.UserID != claims.UserID {
			t.Errorf("Expected UserID %d, got %d", claims.UserID, parsedClaims.UserID)
		}
		if parsedClaims.Username != claims.Username {
			t.Errorf("Expected Username %s, got %s", claims.Username, parsedClaims.Username)
		}
	})

	t.Run("Invalid Token (Mutated)", func(t *testing.T) {
		claims := &service.Claims{UserID: 789}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))

		tamperedToken := tokenString + "mutated"

		_, err := service.ValidateToken(tamperedToken, secret)
		assert.Error(t, err)
	})

	t.Run("Expired Token", func(t *testing.T) {
		timeYesterday := time.Now().Add(-time.Hour * 24)
		claims := &service.Claims{
			UserID: 999,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(timeYesterday),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))

		_, err := service.ValidateToken(tokenString, secret)
		assert.Error(t, err)
	})
}
