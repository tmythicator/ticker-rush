package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)


type Quote struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/api/quote", func(c *gin.Context) {
		fakePrice := 150.0 + (rand.Float64() * 5)

		quote := Quote{
			Symbol:    "AAPL",
			Price:     float64(int(fakePrice*100)) / 100,
			Timestamp: time.Now().Unix(),
		}
		log.Printf("Generated: %s at $%.2f (Struct: %+v)", quote.Symbol, quote.Price, quote)


		c.JSON(http.StatusOK, quote)
	})

	log.Println("Server starting on :8080...")
	r.Run(":8080")
}