// Package apperrors defines application-specific errors.
package apperrors

import "errors"

var (
	// ErrInsufficientFunds is returned when a user does not have enough balance.
	ErrInsufficientFunds = errors.New("insufficient funds")
	// ErrInsufficientQuantity is returned when there are not enough shares to sell.
	ErrInsufficientQuantity = errors.New("insufficient quantity")
	// ErrSymbolNotAllowed is returned when a requested symbol is not in the allowed list.
	ErrSymbolNotAllowed = errors.New("symbol not allowed")
	// ErrAuthRequired is returned when authentication is missing.
	ErrAuthRequired = errors.New("authentication required")
	// ErrInvalidToken is returned when the authentication token is invalid.
	ErrInvalidToken = errors.New("invalid token")
	// ErrMarketClosed is returned when the market is closed.
	ErrMarketClosed = errors.New("market is closed")
)
