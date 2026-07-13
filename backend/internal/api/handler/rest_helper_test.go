package handler_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
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
	ctx           = context.Background()
	sharedConnStr string
)

func TestMain(m *testing.M) {
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
		log.Fatalf("failed to start postgres container: %s", err)
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		if termErr := postgresContainer.Terminate(ctx); termErr != nil {
			log.Printf("failed to terminate postgres container: %s", termErr)
		}
		log.Fatalf("failed to get connection string: %s", err)
	}

	// Run Migrations (Embedded)
	if err := db.Migrate(connStr, "", ""); err != nil {
		if termErr := postgresContainer.Terminate(ctx); termErr != nil {
			log.Printf("failed to terminate postgres container: %s", termErr)
		}
		log.Fatalf("failed to run migrations: %s", err)
	}

	sharedConnStr = connStr

	code := m.Run()

	if termErr := postgresContainer.Terminate(ctx); termErr != nil {
		log.Printf("failed to terminate postgres container: %s", termErr)
	}
	os.Exit(code)
}

type testEnv struct {
	Router             *api.Router
	MiniRedis          *miniredis.Miniredis
	DB                 *pgxpool.Pool
	ValkeyClient       *redis.Client
	LadderRepo         *postgreRepo.LadderRepository
	UserRepo           *postgreRepo.User
	PortfolioRepo      *postgreRepo.PortfolioRepository
	MarketRepo         *redisRepo.MarketRepository
	UserService        *service.User
	TradeService       *service.Trade
	MarketService      *service.Market
	LeaderboardService *service.Leaderboard
}

// MockHistoryRepository mocks the history storage.
type MockHistoryRepository struct{}

func (m *MockHistoryRepository) SaveQuote(ctx context.Context, quote *domain.Quote) error {
	return nil
}

func (m *MockHistoryRepository) GetHistory(ctx context.Context, symbol string, limit int) ([]*domain.Quote, error) {
	return nil, nil
}

func setupTestEnv(t *testing.T) *testEnv {
	mr, _ := miniredis.Run()
	valkeyClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	dbPool, err := pgxpool.New(ctx, sharedConnStr)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	// Truncate all tables to isolate test data
	truncateTables(t, dbPool)

	// Initialize Layers
	ladderRepo := postgreRepo.NewLadderRepository(dbPool)
	userRepo := postgreRepo.NewUser(dbPool)
	portfolioRepo := postgreRepo.NewPortfolioRepository(dbPool)
	marketRepo := redisRepo.NewMarketRepository(valkeyClient)
	leaderboardRepo := redisRepo.NewLeaderboardRepository(valkeyClient)
	transactor := postgreRepo.NewPgxTransactor(dbPool)
	historyRepo := &MockHistoryRepository{}
	rlRepo := redisRepo.NewRateLimitter(valkeyClient)

	userService := service.NewUser(userRepo, portfolioRepo, ladderRepo)
	tradeService := service.NewTrade(userRepo, portfolioRepo, marketRepo, ladderRepo, transactor)
	ladderService := service.NewLadder(ladderRepo)
	leaderboardService := service.NewLeaderboard(userRepo, portfolioRepo, marketRepo, ladderRepo, leaderboardRepo)
	marketService := service.NewMarket(marketRepo, historyRepo, ladderRepo)

	cfg := &config.Config{
		ServerPort: 8080,
		ClientPort: 3000,
		JWTSecret:  testSecret,
	}

	restHandler := handler.NewRestHandler(userService, tradeService, marketService, leaderboardService, ladderService, testSecret)

	router, err := api.NewRouter(restHandler, cfg, rlRepo)
	if err != nil {
		t.Fatalf("Failed to create router: %v", err)
	}

	return &testEnv{
		Router:             router,
		MiniRedis:          mr,
		DB:                 dbPool,
		ValkeyClient:       valkeyClient,
		LadderRepo:         ladderRepo,
		UserRepo:           userRepo,
		PortfolioRepo:      portfolioRepo,
		MarketRepo:         marketRepo,
		UserService:        userService,
		TradeService:       tradeService,
		MarketService:      marketService,
		LeaderboardService: leaderboardService,
	}
}

func truncateTables(t *testing.T, db *pgxpool.Pool) {
	tables := []string{"users"}
	for _, table := range tables {
		_, err := db.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			t.Fatalf("failed to truncate table %s: %v", table, err)
		}
	}
}

// Helper to setup a user and join them to a ladder with a specific balance
func (env *testEnv) setupJoinedUser(t *testing.T, balance float64) (*domain.User, string, int64) {
	createdUser, err := env.UserRepo.CreateUser(ctx, service.CreateUserParams{
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

	u, err := env.UserRepo.GetUser(ctx, createdUser.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	token, _ := service.GenerateToken(u, testSecret)

	wJoin := httptest.NewRecorder()
	reqJoin, _ := http.NewRequest(http.MethodPost, "/api/v1/ladder/participants", nil)
	reqJoin.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	env.Router.ServeHTTP(wJoin, reqJoin)
	if wJoin.Code != http.StatusOK {
		t.Fatalf("Failed to join ladder: %d %s", wJoin.Code, wJoin.Body.String())
	}

	activeLadderID, _ := env.LadderRepo.GetActiveLadder(ctx)

	if balance != 10000.0 {
		err = env.UserRepo.UpdateUserBalance(ctx, u.ID, activeLadderID, decimal.NewFromFloat(balance))
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
