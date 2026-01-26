package service_test

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	pb "github.com/tmythicator/ticker-rush/server/internal/proto/user"
	"github.com/tmythicator/ticker-rush/server/internal/service"
)

func TestAuthService_GenerateToken(t *testing.T) {
	const (
		testEmail  = "test@example.com"
		testUserID = int64(123)
		secret     = "test-secret"
	)

	user := &pb.User{
		Id:    testUserID,
		Email: testEmail,
	}

	err := os.Setenv("JWT_SECRET", secret)
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = os.Unsetenv("JWT_SECRET") }()

	tokenString, err := service.GenerateToken(user)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, float64(testUserID), claims["user_id"])
	assert.Equal(t, testEmail, claims["email"])
}

func TestAuthService_ValidateToken(t *testing.T) {
	const secret = "test-secret"

	err := os.Setenv("JWT_SECRET", secret)
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = os.Unsetenv("JWT_SECRET") }()

	t.Run("Valid Token", func(t *testing.T) {
		claims := &service.Claims{
			UserID: 456,
			Email:  "valid@example.com",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				Issuer:    "ticker-rush",
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))

		parsedClaims, err := service.ValidateToken(tokenString)

		assert.NoError(t, err)
		assert.Equal(t, int64(456), parsedClaims.UserID)
		assert.Equal(t, "valid@example.com", parsedClaims.Email)
	})

	t.Run("Invalid Token (Mutated)", func(t *testing.T) {
		claims := &service.Claims{UserID: 789}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))

		tamperedToken := tokenString + "mutated"

		_, err := service.ValidateToken(tamperedToken)
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

		_, err := service.ValidateToken(tokenString)
		assert.Error(t, err)
	})
}
