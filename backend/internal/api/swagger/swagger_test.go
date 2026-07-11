package swagger

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterRoutes(t *testing.T) {
	// Set Gin to test mode to avoid logging output during test runs
	gin.SetMode(gin.TestMode)

	// Create test router
	r := gin.New()
	api := r.Group("/api")
	RegisterRoutes(api)

	// 1. Test Swagger UI HTML endpoint
	t.Run("Swagger UI HTML page", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/swagger", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
		assert.Contains(t, w.Body.String(), "Ticker Rush API Docs")
		assert.Contains(t, w.Body.String(), "swagger-ui")
	})

	// 2. Test Swagger JSON documents endpoint
	t.Run("Exchange Swagger JSON", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/swagger/docs/exchange/v1/exchange.swagger.json", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

		// Verify it is a valid JSON document
		var doc map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &doc)
		assert.NoError(t, err)

		// Verify basic OpenAPI fields
		assert.Equal(t, "2.0", doc["swagger"])
		info, ok := doc["info"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "Exchange Service API", info["title"])
	})

	t.Run("Leaderboard Swagger JSON", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/swagger/docs/leaderboard/v1/leaderboard.swagger.json", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

		// Verify it is a valid JSON document
		var doc map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &doc)
		assert.NoError(t, err)

		// Verify basic OpenAPI fields
		assert.Equal(t, "2.0", doc["swagger"])
		info, ok := doc["info"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "Leaderboard Service API", info["title"])
	})

	t.Run("User Swagger JSON", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/swagger/docs/user/v1/user.swagger.json", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

		var doc map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &doc)
		assert.NoError(t, err)

		assert.Equal(t, "2.0", doc["swagger"])
		info, ok := doc["info"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "User Service API", info["title"])
	})

	t.Run("Ladder Swagger JSON", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/swagger/docs/ladder/v1/ladder.swagger.json", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

		var doc map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &doc)
		assert.NoError(t, err)

		assert.Equal(t, "2.0", doc["swagger"])
		info, ok := doc["info"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "Ladder Service API", info["title"])
	})

	t.Run("Non-existent Swagger JSON", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/swagger/docs/non-existent.json", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
