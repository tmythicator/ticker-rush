// Package apperrors defines application-specific errors.
package apperrors

import (
	"errors"
	"net/http"
)

// InvalidParam represents a single field validation error.
type InvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

// ProblemDetails represents an RFC 7807 problem details response.
type ProblemDetails struct {
	Type          string         `json:"type"`
	Title         string         `json:"title"`
	Status        int            `json:"status"`
	Detail        string         `json:"detail"`
	Instance      string         `json:"instance"`
	InvalidParams []InvalidParam `json:"invalid_params,omitempty"`
}

// Error type URIs
const (
	TypePrefix            = "https://tickerrush.com/errors/"
	TypeValidation        = TypePrefix + "validation-error"
	TypeInsufficientFunds = TypePrefix + "insufficient-funds"
	TypeInsufficientQty   = TypePrefix + "insufficient-quantity"
	TypeMarketClosed      = TypePrefix + "market-closed"
	TypeAuthRequired      = TypePrefix + "auth-required"
	TypeForbidden         = TypePrefix + "forbidden"
	TypeNotFound          = TypePrefix + "not-found"
	TypeInternalError     = TypePrefix + "internal-error"
	TypeRateLimitExceeded = TypePrefix + "rate-limit-exceeded"
)

// MappedTitle returns a human-readable title for standard problem types.
func MappedTitle(errType string) string {
	switch errType {
	case TypeValidation:
		return "Validation Failed"
	case TypeInsufficientFunds:
		return "Insufficient Funds"
	case TypeInsufficientQty:
		return "Insufficient Stock Quantity"
	case TypeMarketClosed:
		return "Market Closed"
	case TypeAuthRequired:
		return "Authentication Required"
	case TypeForbidden:
		return "Access Forbidden"
	case TypeNotFound:
		return "Resource Not Found"
	case TypeRateLimitExceeded:
		return "Rate Limit Exceeded"
	default:
		return "API Error"
	}
}

// MatchError returns the HTTP status, problem type, and detail message for a given error.
func MatchError(err error) (int, string, string) {
	if err == nil {
		return http.StatusOK, "", ""
	}
	switch {
	case errors.Is(err, ErrAGBNotAccepted),
		errors.Is(err, ErrPasswordTooShort),
		errors.Is(err, ErrPasswordTooLong),
		errors.Is(err, ErrInvalidUsernameFormat),
		errors.Is(err, ErrNameRequired),
		errors.Is(err, ErrProfanityDetected),
		errors.Is(err, ErrUsernameNotAllowed),
		errors.Is(err, ErrInvalidWebsiteFormat),
		errors.Is(err, ErrInvalidQuantity),
		errors.Is(err, ErrAlreadyJoinedLadder):
		return http.StatusBadRequest, TypeValidation, err.Error()

	case errors.Is(err, ErrAuthRequired),
		errors.Is(err, ErrInvalidToken):
		return http.StatusUnauthorized, TypeAuthRequired, err.Error()

	case errors.Is(err, ErrNotJoinedLadder),
		errors.Is(err, ErrLadderNotActive):
		return http.StatusForbidden, TypeForbidden, err.Error()

	case errors.Is(err, ErrPublicProfileNotFoundOrPrivate),
		errors.Is(err, ErrSymbolNotAllowed):
		return http.StatusNotFound, TypeNotFound, err.Error()

	case errors.Is(err, ErrInsufficientFunds):
		return http.StatusPaymentRequired, TypeInsufficientFunds, err.Error()

	case errors.Is(err, ErrInsufficientQuantity):
		return http.StatusBadRequest, TypeInsufficientQty, err.Error()

	case errors.Is(err, ErrMarketClosed):
		return http.StatusForbidden, TypeMarketClosed, err.Error()

	default:
		return http.StatusInternalServerError, TypeInternalError, err.Error()
	}
}

// ValidationErrorParams maps a domain validation error to one or more InvalidParams.
func ValidationErrorParams(err error) []InvalidParam {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, ErrAGBNotAccepted):
		return []InvalidParam{{Name: "agb_accepted", Reason: err.Error()}}
	case errors.Is(err, ErrPasswordTooShort), errors.Is(err, ErrPasswordTooLong):
		return []InvalidParam{{Name: "password", Reason: err.Error()}}
	case errors.Is(err, ErrInvalidUsernameFormat), errors.Is(err, ErrUsernameNotAllowed):
		return []InvalidParam{{Name: "username", Reason: err.Error()}}
	case errors.Is(err, ErrNameRequired):
		return []InvalidParam{
			{Name: "first_name", Reason: err.Error()},
			{Name: "last_name", Reason: err.Error()},
		}
	case errors.Is(err, ErrProfanityDetected):
		// Profanity check can affect username or names
		return []InvalidParam{
			{Name: "username", Reason: err.Error()},
			{Name: "first_name", Reason: err.Error()},
			{Name: "last_name", Reason: err.Error()},
		}
	case errors.Is(err, ErrInvalidWebsiteFormat):
		return []InvalidParam{{Name: "website", Reason: err.Error()}}
	case errors.Is(err, ErrInvalidQuantity):
		return []InvalidParam{{Name: "quantity", Reason: err.Error()}}
	default:
		return nil
	}
}
