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
	// ErrPasswordTooLong is returned when the password is longer than the bcrypt limit (72 bytes).
	ErrPasswordTooLong = errors.New("password must be at most 72 characters long")
	// ErrNameRequired is returned when the first or last name is empty.
	ErrNameRequired = errors.New("first name and last name are required")
	// ErrProfanityDetected is returned when profanity is found in user input.
	ErrProfanityDetected = errors.New("profanity detected in username or name")
	// ErrUsernameNotAllowed is returned when a username is in the blocklist.
	ErrUsernameNotAllowed = errors.New("username is not allowed")
	// ErrAGBNotAccepted is returned when the user has not accepted the terms and conditions.
	ErrAGBNotAccepted = errors.New("you must accept the terms and conditions")
	// ErrNotJoinedLadder is returned when a user tries to trade without joining the ladder.
	ErrNotJoinedLadder = errors.New("user has not joined the active ladder")
	// ErrAlreadyJoinedLadder is returned when a user tries to join a ladder they are already in.
	ErrAlreadyJoinedLadder = errors.New("user has already joined the ladder")
	// ErrLadderNotActive is returned when trading occurs outside the ladder timeframe.
	ErrLadderNotActive = errors.New("ladder is not currently active")
	// ErrInvalidWebsiteFormat is returned when the website format is invalid.
	ErrInvalidWebsiteFormat = errors.New("website must be a valid URL starting with http:// or https://")
	// ErrInvalidQuantity is returned when a trade quantity is invalid.
	ErrInvalidQuantity = errors.New("quantity must be between 0.00000001 and 1,000,000,000")
	// ErrPublicProfileNotFoundOrPrivate is returned when a public profile is requested but not found or is private.
	ErrPublicProfileNotFoundOrPrivate = errors.New("user not found or profile is private")

	// ErrInvalidRequestBody is returned when JSON binding fails.
	ErrInvalidRequestBody = errors.New("invalid request body")
	// ErrFailedToFetchProfileAfterCreation is returned when profile fetch fails after creation.
	ErrFailedToFetchProfileAfterCreation = errors.New("failed to fetch user profile after creation")
	// ErrInvalidCredentials is returned when login credentials do not match.
	ErrInvalidCredentials = errors.New("invalid username or password")
	// ErrFailedToFetchProfile is returned when profile fetch fails.
	ErrFailedToFetchProfile = errors.New("failed to fetch user profile")
	// ErrFailedToGenerateToken is returned when JWT generation fails.
	ErrFailedToGenerateToken = errors.New("failed to generate token")
	// ErrUsernameRequired is returned when username is missing from a request.
	ErrUsernameRequired = errors.New("username is required")
	// ErrInternalServiceError is returned for internal errors.
	ErrInternalServiceError = errors.New("internal service error")
	// ErrSymbolRequired is returned when the ticker symbol is missing.
	ErrSymbolRequired = errors.New("symbol is required")
	// ErrInvalidTradeAction is returned when trade action is not buy or sell.
	ErrInvalidTradeAction = errors.New("invalid trade action")
	// ErrFailedToFetchLeaderboard is returned when leaderboard fetch fails.
	ErrFailedToFetchLeaderboard = errors.New("failed to fetch leaderboard")
	// ErrFailedToFetchActiveLadder is returned when active ladder fetch fails.
	ErrFailedToFetchActiveLadder = errors.New("failed to fetch active ladder")
	// ErrMarketDataWarmingUp is returned when Redis has no quote cache yet.
	ErrMarketDataWarmingUp = errors.New("market data warming up, please retry")
	// ErrInternalAuthConfigurationError is returned when user ID context missing.
	ErrInternalAuthConfigurationError = errors.New("internal authentication configuration error")
)
