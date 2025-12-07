package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tmythicator/ticker-rush/server/internal/service"
)

// TODO: Implement Handler struct and methods here
// It should have methods like GetQuote, BuyStock, SellStock, etc.
// It should inject the Services and handle HTTP requests/responses.

type RestHandler struct {
	userService  *service.UserService
	tradeService *service.TradeService
}

func NewRestHandler(userService *service.UserService, tradeService *service.TradeService) *RestHandler {
	return &RestHandler{
		userService:  userService,
		tradeService: tradeService,
	}
}

type CreateUserRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (h *RestHandler) CreateUser(c *gin.Context) {
	req := CreateUserRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
