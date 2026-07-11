package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/tmythicator/ticker-rush/backend/internal/apperrors"
)

// RespondWithProblem writes an RFC 7807 problem details response to the client.
func RespondWithProblem(
	c *gin.Context,
	status int,
	errType string,
	detail string,
	invalidParams []apperrors.InvalidParam,
) {
	c.Header("Content-Type", "application/problem+json")
	c.JSON(status, apperrors.ProblemDetails{
		Type:          errType,
		Title:         apperrors.MappedTitle(errType),
		Status:        status,
		Detail:        detail,
		Instance:      c.Request.URL.Path,
		InvalidParams: invalidParams,
	})
}
