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

type Router struct {
	engine *gin.Engine
}

func NewRouter(handler *handler.RestHandler, cfg *config.Config) (*Router, error) {
	engine := gin.Default()

	if err := engine.SetTrustedProxies(nil); err != nil {
		return nil, fmt.Errorf("failed to set trusted proxies: %v", err)
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
		api.POST("/register", handler.CreateUser)

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/quote", handler.GetQuote)
			protected.GET("/user/me", handler.GetMe)
			protected.POST("/buy", handler.BuyStock)
			protected.POST("/sell", handler.SellStock)
		}
	}

	return &Router{engine: engine}, nil
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.engine.ServeHTTP(w, req)
}
