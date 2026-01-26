// Package handler provides HTTP handlers.
package handler

// CreateUserRequest represents the payload for creating a new user.
type CreateUserRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// LoginRequest represents the payload for user login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TradeRequest represents the payload for buying or selling stocks.
type TradeRequest struct {
	Symbol string  `json:"symbol"`
	Count  float64 `json:"count"`
}
