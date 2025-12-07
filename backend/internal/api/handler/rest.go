package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	go_redis "github.com/redis/go-redis/v9"
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

type CreateUserRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
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

	c.JSON(http.StatusOK, user)
}

func (h *RestHandler) GetUser(c *gin.Context) {
	idInt, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	user, err := h.userService.GetUserWithPortfolio(c.Request.Context(), idInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Service Error"})
		return
	}
	c.JSON(http.StatusOK, user)
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
	var req model.TradeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.tradeService.BuyStock(c.Request.Context(), req.UserID, req.Symbol, int32(req.Count))
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

	c.JSON(http.StatusOK, user)

}

func (h *RestHandler) SellStock(c *gin.Context) {
	var req model.TradeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.tradeService.SellStock(c.Request.Context(), req.UserID, req.Symbol, int32(req.Count))
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

	c.JSON(http.StatusOK, user)
}
