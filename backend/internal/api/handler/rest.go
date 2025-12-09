package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	go_redis "github.com/redis/go-redis/v9"
	"github.com/tmythicator/ticker-rush/server/internal/api/middleware"
	"github.com/tmythicator/ticker-rush/server/internal/service"
	"github.com/tmythicator/ticker-rush/server/model"
)

type RestHandler struct {
	userService   *service.UserService
	tradeService  *service.TradeService
	marketService *service.MarketService
}

func NewRestHandler(userService *service.UserService, tradeService *service.TradeService, marketService *service.MarketService) *RestHandler {
	return &RestHandler{
		userService:   userService,
		tradeService:  tradeService,
		marketService: marketService,
	}
}

func (h *RestHandler) CreateUser(c *gin.Context) {
	req := CreateUserRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Service Error"})
		return
	}

	c.JSON(http.StatusOK, toUserResponse(user, nil))
}

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

	fullUser, err := h.userService.GetUserWithPortfolio(c.Request.Context(), user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}

	token, err := service.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": toUserResponse(fullUser.User, fullUser.Portfolio)})
}

func (h *RestHandler) GetMe(c *gin.Context) {
	userId, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.userService.GetUserWithPortfolio(c.Request.Context(), userId.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Service Error"})
		return
	}
	c.JSON(http.StatusOK, toUserResponse(user.User, user.Portfolio))
}

func (h *RestHandler) GetQuote(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "AAPL")

	quote, err := h.marketService.GetQuote(c.Request.Context(), symbol)
	if err != nil {
		if errors.Is(err, model.ErrSymbolNotAllowed) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, go_redis.Nil) {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Market data warming up, please retry"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Service Error"})
		return
	}

	c.JSON(http.StatusOK, quote)
}

func (h *RestHandler) BuyStock(c *gin.Context) {
	userId, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req TradeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := h.tradeService.BuyStock(c.Request.Context(), userId.(int64), req.Symbol, req.Count)
	if err != nil {
		if errors.Is(err, model.ErrInsufficientFunds) {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, model.ErrInsufficientQuantity) {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Service Error"})
		return
	}

	fullUser, err := h.userService.GetUserWithPortfolio(c.Request.Context(), userId.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Service Error"})
		return
	}

	c.JSON(http.StatusOK, toUserResponse(fullUser.User, fullUser.Portfolio))

}

func (h *RestHandler) SellStock(c *gin.Context) {
	userId, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req TradeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := h.tradeService.SellStock(c.Request.Context(), userId.(int64), req.Symbol, req.Count)
	if err != nil {
		if errors.Is(err, model.ErrInsufficientQuantity) {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, model.ErrInsufficientFunds) {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Service Error"})
		return
	}

	fullUser, err := h.userService.GetUserWithPortfolio(c.Request.Context(), userId.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Service Error"})
		return
	}

	c.JSON(http.StatusOK, toUserResponse(fullUser.User, fullUser.Portfolio))
}
