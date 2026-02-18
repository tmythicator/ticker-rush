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
	// ErrInvalidUsernameFormat is returned when the username does not match the required format.
	ErrInvalidUsernameFormat = errors.New("invalid username: must be 3-20 alphanumeric characters or underscores")
	// ErrPasswordTooShort is returned when the password is shorter than the required length.
	ErrPasswordTooShort = errors.New("password must be at least 8 characters long")
	// ErrNameRequired is returned when the first or last name is empty.
	ErrNameRequired = errors.New("first name and last name are required")
	// ErrProfanityDetected is returned when profanity is found in user input.
	ErrProfanityDetected = errors.New("profanity detected in username or name")
	// ErrUserNotFound is returned when a user is not found.
	ErrUserNotFound = errors.New("user not found")
	// ErrUsernameNotAllowed is returned when a username is in the blocklist.
	ErrUsernameNotAllowed = errors.New("username is not allowed")
)
