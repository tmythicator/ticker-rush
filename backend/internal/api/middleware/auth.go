// Package middleware provides HTTP middlewares.
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tmythicator/ticker-rush/backend/internal/apperrors"
	"github.com/tmythicator/ticker-rush/backend/internal/service"
)

// UserIDKey is the key used to store the user ID in the Gin context.
const UserIDKey = "userID"

// AuthMiddleware is a Gin middleware that validates the authentication token from cookies.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("auth_token")
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, apperrors.ErrAuthRequired)

			return
		}

		if token == "" {
			_ = c.AbortWithError(http.StatusUnauthorized, apperrors.ErrInvalidToken)

			return
		}

		claims, err := service.ValidateToken(token, jwtSecret)
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, apperrors.ErrInvalidToken)

			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}

// GetUserID retrieves the authenticated user's ID from the context.
func GetUserID(c *gin.Context) (int64, bool) {
	val, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}
	id, ok := val.(int64)
	return id, ok
}
