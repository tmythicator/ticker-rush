// Package handler provides the HTTP handlers for the API.
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

	"github.com/tmythicator/ticker-rush/backend/internal/api/middleware"
	"github.com/tmythicator/ticker-rush/backend/internal/apperrors"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/exchange/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/ladder/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/user/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/service"
)

// RestHandler handles HTTP requests for the API.
type RestHandler struct {
	userService   *service.User
	tradeService  *service.Trade
	marketService *service.Market
	leadService   *service.Leaderboard
	ladderService *service.Ladder
	jwtSecret     string
}

// NewRestHandler creates a new instance of RestHandler.
func NewRestHandler(
	userService *service.User,
	tradeService *service.Trade,
	marketService *service.Market,
	leadService *service.Leaderboard,
	ladderService *service.Ladder,
	jwtSecret string,
) *RestHandler {
	return &RestHandler{
		userService:   userService,
		tradeService:  tradeService,
		marketService: marketService,
		leadService:   leadService,
		ladderService: ladderService,
		jwtSecret:     jwtSecret,
	}
}

// CreateUser handles user registration.
func (h *RestHandler) CreateUser(c *gin.Context) {
	req := user.CreateUserRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithProblem(c, http.StatusBadRequest, apperrors.TypeValidation, "Invalid request body", nil)

		return
	}

	createdUser, err := h.userService.CreateUser(
		c.Request.Context(),
		req.Username,
		req.Password,
		req.FirstName,
		req.LastName,
		req.AgbAccepted,
	)
	if err != nil {
		status, errType, detail := apperrors.MatchError(err)
		invalidParams := apperrors.ValidationErrorParams(err)
		RespondWithProblem(c, status, errType, detail, invalidParams)

		return
	}

	fullUser, err := h.userService.GetUserWithPortfolio(c.Request.Context(), createdUser.ID)
	if err != nil {
		RespondWithProblem(c, http.StatusInternalServerError, apperrors.TypeInternalError, "Failed to fetch user profile after creation", nil)

		return
	}

	c.JSON(http.StatusOK, &user.CreateUserResponse{User: ToExternalUser(fullUser)})
}

// Login handles user authentication.
func (h *RestHandler) Login(c *gin.Context) {
	var req user.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithProblem(c, http.StatusBadRequest, apperrors.TypeValidation, "Invalid request body", nil)

		return
	}

	authenticatedUser, err := h.userService.Authenticate(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		RespondWithProblem(c, http.StatusUnauthorized, apperrors.TypeAuthRequired, "Invalid username or password", nil)

		return
	}

	fullUser, err := h.userService.GetUserWithPortfolio(c.Request.Context(), authenticatedUser.ID)
	if err != nil {
		RespondWithProblem(c, http.StatusInternalServerError, apperrors.TypeInternalError, "Failed to fetch user profile", nil)

		return
	}

	token, err := service.GenerateToken(authenticatedUser, h.jwtSecret)
	if err != nil {
		RespondWithProblem(c, http.StatusInternalServerError, apperrors.TypeInternalError, "Failed to generate token", nil)

		return
	}

	c.SetCookie("auth_token", token, int(service.SessionDuration.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, &user.LoginResponse{User: ToExternalUser(fullUser)})
}

// Logout handles user logout.
func (h *RestHandler) Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// UpdateUser handles user profile updates.
func (h *RestHandler) UpdateUser(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	var req user.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithProblem(c, http.StatusBadRequest, apperrors.TypeValidation, "Invalid request body", nil)

		return
	}

	updatedUser, err := h.userService.UpdateUser(
		c.Request.Context(),
		userID,
		req.FirstName,
		req.LastName,
		req.Website,
		req.IsPublic,
	)

	if err != nil {
		status, errType, detail := apperrors.MatchError(err)
		invalidParams := apperrors.ValidationErrorParams(err)
		RespondWithProblem(c, status, errType, detail, invalidParams)

		return
	}

	c.JSON(http.StatusOK, &user.UpdateUserResponse{User: ToExternalUser(updatedUser)})
}

// GetPublicProfile handles retrieving a user's public profile.
func (h *RestHandler) GetPublicProfile(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		RespondWithProblem(c, http.StatusBadRequest, apperrors.TypeValidation, "Username is required", nil)

		return
	}

	publicProfile, err := h.userService.GetPublicProfile(c.Request.Context(), username)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			RespondWithProblem(c, http.StatusNotFound, apperrors.TypeNotFound, "User not found or profile is private", nil)

			return
		}
		status, errType, detail := apperrors.MatchError(err)
		invalidParams := apperrors.ValidationErrorParams(err)
		RespondWithProblem(c, status, errType, detail, invalidParams)

		return
	}

	c.JSON(http.StatusOK, &user.GetPublicProfileResponse{User: ToExternalUser(publicProfile)})
}

// GetMe returns the current user's profile.
func (h *RestHandler) GetMe(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	fullUser, err := h.userService.GetUserWithPortfolio(c.Request.Context(), userID)
	if err != nil {
		RespondWithProblem(c, http.StatusInternalServerError, apperrors.TypeInternalError, "Internal Service Error", nil)

		return
	}

	c.JSON(http.StatusOK, &user.GetMeResponse{User: ToExternalUser(fullUser)})
}

// GetQuote returns a stock quote for a given symbol.
func (h *RestHandler) GetQuote(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		RespondWithProblem(c, http.StatusBadRequest, apperrors.TypeValidation, "Symbol is required", nil)

		return
	}

	quote, err := h.marketService.GetQuote(c.Request.Context(), symbol)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			RespondWithProblem(
				c,
				http.StatusServiceUnavailable,
				apperrors.TypeInternalError,
				"Market data warming up, please retry",
				nil,
			)

			return
		}

		status, errType, detail := apperrors.MatchError(err)
		invalidParams := apperrors.ValidationErrorParams(err)
		RespondWithProblem(c, status, errType, detail, invalidParams)

		return
	}

	c.JSON(http.StatusOK, &exchange.GetQuoteResponse{Quote: ToExternalQuote(quote)})
}

// CreateTrade handles stock buy or sell requests.
func (h *RestHandler) CreateTrade(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	var req exchange.CreateTradeRequest
	if err := c.BindJSON(&req); err != nil {
		RespondWithProblem(c, http.StatusBadRequest, apperrors.TypeValidation, "Invalid request body", nil)

		return
	}

	var err error
	var msg string
	switch req.Action {
	case exchange.TradeAction_TRADE_ACTION_BUY:
		_, err = h.tradeService.BuyStock(c.Request.Context(), userID, req.Symbol, req.Quantity)
		msg = "Bought successfully"
	case exchange.TradeAction_TRADE_ACTION_SELL:
		_, err = h.tradeService.SellStock(c.Request.Context(), userID, req.Symbol, req.Quantity)
		msg = "Sold successfully"
	default:
		RespondWithProblem(c, http.StatusBadRequest, apperrors.TypeValidation, "Invalid trade action", nil)

		return
	}

	if err != nil {
		status, errType, detail := apperrors.MatchError(err)
		invalidParams := apperrors.ValidationErrorParams(err)
		RespondWithProblem(c, status, errType, detail, invalidParams)

		return
	}

	fullUser, err := h.userService.GetUserWithPortfolio(c.Request.Context(), userID)
	if err != nil {
		RespondWithProblem(c, http.StatusInternalServerError, apperrors.TypeInternalError, "Internal Service Error", nil)

		return
	}

	c.JSON(http.StatusOK, &exchange.CreateTradeResponse{
		Success: true,
		Message: msg,
		Participant: &ladder.LadderParticipant{
			User: ToExternalUser(fullUser),
		},
	})
}

// GetLeaderboard handles leaderboard fetching requests.
func (h *RestHandler) GetLeaderboard(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	resp, err := h.leadService.GetLeaderboard(c.Request.Context(), offset, limit)
	if err != nil {
		RespondWithProblem(c, http.StatusInternalServerError, apperrors.TypeInternalError, "Failed to fetch leaderboard", nil)

		return
	}
	c.JSON(http.StatusOK, ToExternalLeaderboardResponse(resp))
}

// GetActiveLadder returns full metadata for the currently active ladder.
func (h *RestHandler) GetActiveLadder(c *gin.Context) {
	l, err := h.ladderService.GetActiveLadder(c.Request.Context())
	if err != nil {
		RespondWithProblem(c, http.StatusInternalServerError, apperrors.TypeInternalError, "Failed to fetch active ladder", nil)

		return
	}

	c.JSON(http.StatusOK, &ladder.GetActiveLadderResponse{
		Ladder: ToExternalLadder(l),
	})
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
		RespondWithProblem(c, http.StatusBadRequest, apperrors.TypeValidation, err.Error(), nil)

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

// GetHistory handles requests for historical market data.
func (h *RestHandler) GetHistory(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		RespondWithProblem(c, http.StatusBadRequest, apperrors.TypeValidation, "symbol is required", nil)

		return
	}

	limit := 100
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = val
		}
	}

	history, err := h.marketService.GetHistory(c.Request.Context(), symbol, limit)
	if err != nil {
		status, errType, detail := apperrors.MatchError(err)
		invalidParams := apperrors.ValidationErrorParams(err)
		RespondWithProblem(c, status, errType, detail, invalidParams)

		return
	}

	protoHistory := make([]*exchange.Quote, len(history))
	for i, q := range history {
		protoHistory[i] = ToExternalQuote(q)
	}

	c.JSON(http.StatusOK, &exchange.GetHistoryResponse{
		History: protoHistory,
	})
}

// JoinLadder allows a user to join the active ladder.
func (h *RestHandler) JoinLadder(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	err := h.ladderService.JoinLadder(c.Request.Context(), userID)
	if err != nil {
		status, errType, detail := apperrors.MatchError(err)
		invalidParams := apperrors.ValidationErrorParams(err)
		RespondWithProblem(c, status, errType, detail, invalidParams)

		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Joined ladder successfully"})
}

// GetUserID retrieves the authenticated user ID from context and aborts with a 500 status if missing.
func (h *RestHandler) getUserID(c *gin.Context) (int64, bool) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		RespondWithProblem(
			c,
			http.StatusInternalServerError,
			apperrors.TypeInternalError,
			"Internal authentication configuration error",
			nil,
		)

		return 0, false
	}

	return userID, true
}
