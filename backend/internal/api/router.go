// Package api provides the API router and setup.
package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tmythicator/ticker-rush/server/internal/api/handler"
	"github.com/tmythicator/ticker-rush/server/internal/api/middleware"
	"github.com/tmythicator/ticker-rush/server/internal/config"
)

// Router handles API routing.
type Router struct {
	engine *gin.Engine
}

// NewRouter creates a new API router.
func NewRouter(handler *handler.RestHandler, cfg *config.Config) (*Router, error) {
	engine := gin.Default()

	err := engine.SetTrustedProxies(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to set trusted proxies: %w", err)
	}

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{fmt.Sprintf("http://localhost:%d", cfg.ClientPort)},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := engine.Group("/api")
	{
		api.POST("/login", handler.Login)
		api.POST("/logout", handler.Logout)
		api.POST("/register", handler.CreateUser)
		api.GET("/config", handler.GetConfig)

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			protected.GET("/quotes/events", handler.StreamQuotes)
			protected.GET("/quote", handler.GetQuote)
			protected.GET("/history", handler.GetHistory)
			protected.GET("/user/me", handler.GetMe)
			protected.POST("/buy", handler.BuyStock)
			protected.POST("/sell", handler.SellStock)
			protected.GET("/leaderboard", handler.GetLeaderboard)
		}
	}

	return &Router{engine: engine}, nil
}

// Run starts the HTTP server.
func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.engine.ServeHTTP(w, req)
}
