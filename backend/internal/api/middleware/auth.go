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
			respondWithProblemDirectly(c, http.StatusUnauthorized, apperrors.TypeAuthRequired, apperrors.ErrAuthRequired.Error())
			c.Abort()

			return
		}

		if token == "" {
			respondWithProblemDirectly(c, http.StatusUnauthorized, apperrors.TypeAuthRequired, apperrors.ErrInvalidToken.Error())
			c.Abort()

			return
		}

		claims, err := service.ValidateToken(token, jwtSecret)
		if err != nil {
			respondWithProblemDirectly(c, http.StatusUnauthorized, apperrors.TypeAuthRequired, apperrors.ErrInvalidToken.Error())
			c.Abort()

			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}

func respondWithProblemDirectly(c *gin.Context, status int, errType string, detail string) {
	c.Header("Content-Type", "application/problem+json")
	c.JSON(status, apperrors.ProblemDetails{
		Type:     errType,
		Title:    apperrors.MappedTitle(errType),
		Status:   status,
		Detail:   detail,
		Instance: c.Request.URL.Path,
	})
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
