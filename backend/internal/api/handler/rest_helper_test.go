package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/tmythicator/ticker-rush/backend/db"
	"github.com/tmythicator/ticker-rush/backend/internal/api"
	"github.com/tmythicator/ticker-rush/backend/internal/api/handler"
	"github.com/tmythicator/ticker-rush/backend/internal/config"
	"github.com/tmythicator/ticker-rush/backend/internal/domain"
	postgreRepo "github.com/tmythicator/ticker-rush/backend/internal/repository/postgres"
	redisRepo "github.com/tmythicator/ticker-rush/backend/internal/repository/redis"
	"github.com/tmythicator/ticker-rush/backend/internal/service"
)

const (
	testUsername = "test_user"
	testSecret   = "test-secret"
)

var (
	ctx                = context.Background()
	valkeyClient       *redis.Client
	dbPool             *pgxpool.Pool
	ladderRepo         *postgreRepo.LadderRepository
	userRepo           *postgreRepo.User
	portfolioRepo      *postgreRepo.PortfolioRepository
	marketRepo         *redisRepo.MarketRepository
	userService        *service.User
	tradeService       *service.Trade
	marketService      *service.Market
	leaderboardService *service.Leaderboard
)

func setupTestPostgres(t *testing.T) string {
	dbName := "test_db"
	dbUser := "test_user"
	dbPassword := "test_password"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %s", err)
	}

	t.Cleanup(func() {
		termErr := postgresContainer.Terminate(ctx)
		if termErr != nil {
			t.Fatalf("failed to terminate postgres container: %s", termErr)
		}
	})

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	// Run Migrations (Embedded)
	if err := db.Migrate(connStr, "", ""); err != nil {
		t.Fatalf("failed to run migrations: %s", err)
	}

	return connStr
}

// MockHistoryRepository mocks the history storage.
type MockHistoryRepository struct{}

func (m *MockHistoryRepository) SaveQuote(ctx context.Context, quote *domain.Quote) error {
	return nil
}

func (m *MockHistoryRepository) GetHistory(ctx context.Context, symbol string, limit int) ([]*domain.Quote, error) {
	return nil, nil // Return empty history for tests
}

func setupTestRouter(t *testing.T) (*api.Router, *miniredis.Miniredis, *pgxpool.Pool) {
	mr, _ := miniredis.Run()
	valkeyClient = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	connStr := setupTestPostgres(t)

	var err error

	dbPool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	// Initialize Layers
	ladderRepo = postgreRepo.NewLadderRepository(dbPool)
	userRepo = postgreRepo.NewUser(dbPool)
	portfolioRepo = postgreRepo.NewPortfolioRepository(dbPool)
	marketRepo = redisRepo.NewMarketRepository(valkeyClient)
	leaderboardRepo := redisRepo.NewLeaderboardRepository(valkeyClient)
	transactor := postgreRepo.NewPgxTransactor(dbPool)
	historyRepo := &MockHistoryRepository{}
	rlRepo := redisRepo.NewRateLimitter(valkeyClient)

	userService = service.NewUser(userRepo, portfolioRepo, ladderRepo)
	tradeService = service.NewTrade(userRepo, portfolioRepo, marketRepo, ladderRepo, transactor)
	ladderService := service.NewLadder(ladderRepo)
	leaderboardService = service.NewLeaderboard(userRepo, portfolioRepo, marketRepo, ladderRepo, leaderboardRepo)

	marketService = service.NewMarket(marketRepo, historyRepo, ladderRepo)

	cfg := &config.Config{
		ServerPort: 8080,
		ClientPort: 3000,
		JWTSecret:  testSecret,
	}

	restHandler := handler.NewRestHandler(userService, tradeService, marketService, leaderboardService, ladderService, testSecret)

	// Initialize Router
	router, err := api.NewRouter(restHandler, cfg, rlRepo)
	if err != nil {
		t.Fatalf("Failed to create router: %v", err)
	}

	return router, mr, dbPool
}

// Helper to setup a user and join them to a ladder with a specific balance
func setupJoinedUser(ctx context.Context, t *testing.T, r *api.Router, balance float64) (*domain.User, string, int64) {
	createdUser, err := userRepo.CreateUser(ctx, service.CreateUserParams{
		Username:      testUsername,
		PasswordHash:  "password123",
		FirstName:     "Test",
		LastName:      "User",
		Website:       "",
		IsPublic:      false,
		AgbAcceptedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	u, err := userRepo.GetUser(ctx, createdUser.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	token, _ := service.GenerateToken(u, testSecret)

	// Join Ladder
	wJoin := httptest.NewRecorder()
	reqJoin, _ := http.NewRequest(http.MethodPost, "/api/v1/ladder/participants", nil)
	reqJoin.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	r.ServeHTTP(wJoin, reqJoin)
	if wJoin.Code != http.StatusOK {
		t.Fatalf("Failed to join ladder: %d %s", wJoin.Code, wJoin.Body.String())
	}

	activeLadderID, _ := ladderRepo.GetActiveLadder(ctx)

	// Override balance if needed
	if balance != 10000.0 {
		err = userRepo.UpdateUserBalance(ctx, u.ID, activeLadderID, decimal.NewFromFloat(balance))
		if err != nil {
			t.Fatalf("Failed to override balance: %v", err)
		}
	}

	return u, token, activeLadderID
}

func assertPublicProfilePrivacy(t *testing.T, userMap map[string]interface{}) {
	assert.Nil(t, userMap["is_admin"], "PublicProfile should not leak is_admin")
	assert.Nil(t, userMap["is_banned"], "PublicProfile should not leak is_banned")
	assert.Nil(t, userMap["created_at"], "PublicProfile should not leak created_at")
	assert.Nil(t, userMap["is_participating"], "PublicProfile should not leak is_participating")
}
