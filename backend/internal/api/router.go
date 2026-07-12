// Package api provides the API router and setup.
package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/tmythicator/ticker-rush/backend/internal/api/handler"
	"github.com/tmythicator/ticker-rush/backend/internal/api/middleware"
	"github.com/tmythicator/ticker-rush/backend/internal/api/swagger"
	"github.com/tmythicator/ticker-rush/backend/internal/config"
)

// Router handles API routing.
type Router struct {
	engine *gin.Engine
}

// NewRouter creates a new API router.
func NewRouter(handler *handler.RestHandler, cfg *config.Config, rateLimitRepo middleware.RateLimit) (*Router, error) {
	engine := gin.Default()

	if cfg.OtelEndpoint != "" {
		engine.Use(otelgin.Middleware(cfg.OtelServiceName))
	}

	globalLimiter := middleware.NewRateLimitter(rateLimitRepo, 100, time.Minute)
	strictLimiter := middleware.NewRateLimitter(rateLimitRepo, 30, time.Minute)

	err := engine.SetTrustedProxies(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to set trusted proxies: %w", err)
	}

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{fmt.Sprintf("http://localhost:%d", cfg.ClientPort)},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset", "Retry-After"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v1 := engine.Group("/api/v1")
	v1.Use(globalLimiter.Limit())
	{
		v1.POST("/sessions", handler.Login)
		v1.DELETE("/sessions", handler.Logout)
		v1.POST("/users", handler.CreateUser)
		v1.GET("/ladder/active", handler.GetActiveLadder)
		v1.GET("/quotes/:symbol/history", handler.GetHistory)
		v1.GET("/leaderboard", handler.GetLeaderboard)
		v1.GET("/users/:username", handler.GetPublicProfile)

		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		protected.Use(strictLimiter.Limit())
		{
			protected.POST("/ladder/participants", handler.JoinLadder)
			protected.GET("/quotes/events", handler.StreamQuotes)
			protected.GET("/quotes/:symbol", handler.GetQuote)
			protected.GET("/profile", handler.GetMe)
			protected.PUT("/profile", handler.UpdateUser)
			protected.POST("/trades", handler.CreateTrade)
		}
	}

	swagger.RegisterRoutes(engine.Group("/api"))

	return &Router{engine: engine}, nil
}

// Run starts the HTTP server.
func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.engine.ServeHTTP(w, req)
}
