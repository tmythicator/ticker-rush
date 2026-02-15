package handler

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/tmythicator/ticker-rush/server/internal/api/middleware"
	"github.com/tmythicator/ticker-rush/server/internal/apperrors"
	"github.com/tmythicator/ticker-rush/server/internal/proto/leaderboard/v1"
	"github.com/tmythicator/ticker-rush/server/internal/service"
)

// RestHandler handles HTTP requests for the API.
type RestHandler struct {
	userService   *service.UserService
	tradeService  *service.TradeService
	marketService *service.MarketService
	leadService   *service.LeaderBoardService
	jwtSecret     string
}

// NewRestHandler creates a new instance of RestHandler.
func NewRestHandler(
	userService *service.UserService,
	tradeService *service.TradeService,
	marketService *service.MarketService,
	leadService *service.LeaderBoardService,
	jwtSecret string,
) *RestHandler {
	return &RestHandler{
		userService:   userService,
		tradeService:  tradeService,
		marketService: marketService,
		leadService:   leadService,
		jwtSecret:     jwtSecret,
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

	token, err := service.GenerateToken(user, h.jwtSecret)
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

		if errors.Is(err, redis.Nil) {
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

// GetLeaderboard handles leaderboard fetching requests.
func (h *RestHandler) GetLeaderboard(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	resp, err := h.leadService.GetLeaderboard(c.Request.Context(), &leaderboard.GetLeaderboardRequest{Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard"})

		return
	}
	c.JSON(http.StatusOK, resp)
}

// StreamQuotes handles SSE connection for real-time quotes.
func (h *RestHandler) StreamQuotes(c *gin.Context) {
	c.Header("X-Accel-Buffering", "no")
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	symbol := c.Query("symbol")
	log.Printf("[StreamQuotes] Starting stream for symbol: %s", symbol)

	clientGone := c.Writer.CloseNotify()

	pubsub, err := h.marketService.SubscribeToQuotes(c.Request.Context(), symbol)
	if err != nil {
		log.Printf("[StreamQuotes] Failed to subscribe: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	defer func() {
		log.Printf("[StreamQuotes] Closing pubsub for symbol: %s", symbol)
		err := pubsub.Close()
		if err != nil {
			log.Printf("Failed to close pubsub: %v", err)
		}
	}()

	ch := pubsub.Channel()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			log.Printf("[StreamQuotes] Client gone regarding symbol: %s", symbol)

			return false
		case <-ticker.C:
			h.SendHeartbeat(w)

			return true
		case msg := <-ch:
			c.SSEvent("quote", msg.Payload)
			// Flush the response to ensure it's sent immediately
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}

			return true
		}
	})
}

// SendHeartbeat keeps the connection alive by sending a comment.
func (h *RestHandler) SendHeartbeat(w io.Writer) {
	_, _ = w.Write([]byte(": keep-alive\n\n"))
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}
