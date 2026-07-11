package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/tmythicator/ticker-rush/backend/db"
	postgresRepo "github.com/tmythicator/ticker-rush/backend/internal/repository/postgres"
	"github.com/tmythicator/ticker-rush/backend/internal/service"
)

func TestLadderRepository_GetLadderParticipants_WithNullFinalBalance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	dbName := "test_db"
	dbUser := "test_user"
	dbPassword := "test_password"

	// Start Postgres container
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
	require.NoError(t, err)

	t.Cleanup(func() {
		termErr := postgresContainer.Terminate(ctx)
		if termErr != nil {
			t.Fatalf("failed to terminate postgres container: %s", termErr)
		}
	})

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	// Run Migrations
	err = db.Migrate(connStr, "", "")
	require.NoError(t, err)

	// Connect pool
	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)
	defer pool.Close()

	// Initialize repositories
	ladderRepo := postgresRepo.NewLadderRepository(pool)
	userRepo := postgresRepo.NewUser(pool)

	// 1. Create a user
	createdUser, err := userRepo.CreateUser(ctx, service.CreateUserParams{
		Username:      "participant_1",
		PasswordHash:  "hash",
		FirstName:     "First",
		LastName:      "Last",
		Website:       "",
		IsPublic:      true,
		AgbAcceptedAt: time.Now(),
	})
	require.NoError(t, err)

	// 2. Create a ladder
	var ladderID int64
	err = pool.QueryRow(ctx, `
		INSERT INTO ladders (name, type, start_time, end_time, initial_balance, is_active)
		VALUES ('Test Season 1', 'monthly', NOW(), NOW() + INTERVAL '30 days', 10000.0, true)
		RETURNING id`).Scan(&ladderID)
	require.NoError(t, err)

	// 3. Join the ladder (this creates ladder_participant with final_balance = NULL)
	err = ladderRepo.JoinLadder(ctx, ladderID, createdUser.ID)
	require.NoError(t, err)

	// 4. Retrieve participants
	participants, err := ladderRepo.GetLadderParticipants(ctx, ladderID)
	require.NoError(t, err)

	// 5. Assertions
	require.Len(t, participants, 1)
	assert.Equal(t, createdUser.ID, participants[0].User.ID)
	assert.True(t, participants[0].FinalBalance.IsZero())
}
