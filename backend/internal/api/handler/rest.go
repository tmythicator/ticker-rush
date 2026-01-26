package handler

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	go_redis "github.com/redis/go-redis/v9"
	"github.com/tmythicator/ticker-rush/server/internal/api/middleware"
	"github.com/tmythicator/ticker-rush/server/internal/apperrors"
	"github.com/tmythicator/ticker-rush/server/internal/service"
)

// RestHandler handles HTTP requests for the API.
type RestHandler struct {
	userService   *service.UserService
	tradeService  *service.TradeService
	marketService *service.MarketService
}

// NewRestHandler creates a new instance of RestHandler.
func NewRestHandler(
	userService *service.UserService,
	tradeService *service.TradeService,
	marketService *service.MarketService,
) *RestHandler {
	return &RestHandler{
		userService:   userService,
		tradeService:  tradeService,
		marketService: marketService,
	}
}

// CreateUser handles user registration.
func (h *RestHandler) CreateUser(c *gin.Context) {
	req := CreateUserRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})

		return
	}

	user, err := h.userService.CreateUser(
		c.Request.Context(),
		req.Email,
		req.Password,
		req.FirstName,
		req.LastName,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Service Error"})

		return
	}

	c.JSON(http.StatusOK, user)
}

// Login handles user authentication.
func (h *RestHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})

		return
	}

	user, err := h.userService.Authenticate(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})

		return
	}

	fullUser, err := h.userService.GetUserWithPortfolio(c.Request.Context(), user.GetId())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})

		return
	}

	token, err := service.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})

		return
	}

	c.SetCookie("auth_token", token, 3600*24, "/", "", false, true)
	c.JSON(http.StatusOK, fullUser)
}

// Logout handles user logout.
func (h *RestHandler) Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// GetMe returns the current user's profile.
func (h *RestHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})

		return
	}

	user, err := h.userService.GetUserWithPortfolio(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Service Error"})

		return
	}

	c.JSON(http.StatusOK, user)
}

// GetQuote returns a stock quote for a given symbol.
func (h *RestHandler) GetQuote(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "AAPL")

	quote, err := h.marketService.GetQuote(c.Request.Context(), symbol)
	if err != nil {
		if errors.Is(err, apperrors.ErrSymbolNotAllowed) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

			return
		}

		if errors.Is(err, go_redis.Nil) {
			c.JSON(
				http.StatusServiceUnavailable,
				gin.H{"error": "Market data warming up, please retry"},
			)

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Service Error"})

		return
	}

	c.JSON(http.StatusOK, quote)
}

// BuyStock handles stock purchase requests.
func (h *RestHandler) BuyStock(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})

		return
	}

	var req TradeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})

		return
	}

	_, err := h.tradeService.BuyStock(c.Request.Context(), userID.(int64), req.Symbol, req.Count)
	if err != nil {
		if errors.Is(err, apperrors.ErrInsufficientFunds) {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})

			return
		}

		if errors.Is(err, apperrors.ErrInsufficientQuantity) {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Service Error"})

		return
	}

	fullUser, err := h.userService.GetUserWithPortfolio(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Service Error"})

		return
	}

	c.JSON(http.StatusOK, fullUser)
}

// SellStock handles stock sale requests.
func (h *RestHandler) SellStock(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})

		return
	}

	var req TradeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})

		return
	}

	_, err := h.tradeService.SellStock(c.Request.Context(), userID.(int64), req.Symbol, req.Count)
	if err != nil {
		if errors.Is(err, apperrors.ErrInsufficientQuantity) {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})

			return
		}

		if errors.Is(err, apperrors.ErrInsufficientFunds) {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Service Error"})

		return
	}

	fullUser, err := h.userService.GetUserWithPortfolio(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Service Error"})

		return
	}

	c.JSON(http.StatusOK, fullUser)
}

// StreamQuotes handles SSE connection for real-time quotes.
func (h *RestHandler) StreamQuotes(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	clientGone := c.Writer.CloseNotify()

	pubsub, err := h.marketService.SubscribeToQuotes(c.Request.Context(), c.Query("symbol"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	defer func() {
		err := pubsub.Close()
		if err != nil {
			log.Printf("Failed to close pubsub: %v", err)
		}
	}()

	ch := pubsub.Channel()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			return false
		case msg := <-ch:
			c.SSEvent("quote", msg.Payload)

			return true
		}
	})
}
