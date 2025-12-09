package handler

import (
	pb "github.com/tmythicator/ticker-rush/server/proto/user"
)

type UserResponse struct {
	Email     string                       `json:"email"`
	FirstName string                       `json:"first_name"`
	LastName  string                       `json:"last_name"`
	Balance   float64                      `json:"balance"`
	Portfolio map[string]*pb.PortfolioItem `json:"portfolio,omitempty"`
}

func toUserResponse(user *pb.User, portfolio map[string]*pb.PortfolioItem) UserResponse {
	return UserResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Balance:   user.Balance,
		Portfolio: portfolio,
	}
}

type CreateUserRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TradeRequest struct {
	Symbol string  `json:"symbol"`
	Count  float64 `json:"count"`
}
