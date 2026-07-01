// Package main seeds the database with mock data.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/tmythicator/ticker-rush/backend/internal/config"
)

func main() {
	ctx := context.Background()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if os.Getenv("APP_ENV") == "production" {
		log.Fatalf("CRITICAL ERROR: Seeding is disabled in production environment!")
	}

	pool, err := setupDB(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}
	defer pool.Close()

	userIDs, err := seedUsers(ctx, pool)
	if err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}

	if err := seedLadders(ctx, pool, userIDs); err != nil {
		log.Fatalf("Failed to seed ladders: %v", err)
	}

	log.Println("Database seeding completed successfully!")
}

// setupDB connects to the database and truncates existing data.
func setupDB(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	connStr := cfg.DatabaseURL()
	log.Printf("Connecting to database for seeding...")
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	log.Println("Truncating tables...")
	_, err = pool.Exec(ctx, "TRUNCATE TABLE users, ladders, ladder_participants, ladder_portfolio_items, ladder_tickers CASCADE")
	if err != nil {
		pool.Close()

		return nil, fmt.Errorf("failed to truncate tables: %w", err)
	}

	return pool, nil
}

// seedUsers generates mock users and returns their database IDs.
func seedUsers(ctx context.Context, pool *pgxpool.Pool) ([]int64, error) {
	log.Println("Generating password hash...")
	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	passwordHash := string(passwordHashBytes)

	log.Println("Inserting mock users...")
	mockUsers := []struct {
		username  string
		firstName string
		lastName  string
		isAdmin   bool
	}{
		{"admin", "System", "Administrator", true},
		{"alex_t", "Ali", "Time", false},
		{"john_doe", "John", "Doe", false},
		{"jane_smith", "Jane", "Smith", false},
		{"trader_bob", "Bob", "Trader", false},
		{"speedy_whale", "Speedy", "Whale", false},
		{"hft_bot_1", "HighFreq", "Bot One", false},
		{"hft_bot_2", "HighFreq", "Bot Two", false},
		{"noob_investor", "Noob", "Investor", false},
		{"lucky_guy", "Lucky", "Guy", false},
	}

	var userIDs []int64
	for _, u := range mockUsers {
		var id int64
		err = pool.QueryRow(ctx, `
			INSERT INTO users (username, password_hash, first_name, last_name, website, created_at, is_public, is_admin, agb_accepted_at)
			VALUES ($1, $2, $3, $4, '', NOW(), true, $5, NOW())
			RETURNING id`,
			u.username, passwordHash, u.firstName, u.lastName, u.isAdmin).Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("failed to insert user %s: %w", u.username, err)
		}
		userIDs = append(userIDs, id)
	}

	return userIDs, nil
}

// seedLadders generates completed, active, and upcoming ladders and registers users.
func seedLadders(ctx context.Context, pool *pgxpool.Pool, userIDs []int64) error {
	now := time.Now().UTC()
	currentYear, currentMonth, _ := now.Date()

	// Calculate dates
	startOfCurrent := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.UTC)
	endOfCurrent := startOfCurrent.AddDate(0, 1, 0).Add(-time.Second)

	startOfPast1 := startOfCurrent.AddDate(0, -1, 0)
	endOfPast1 := startOfCurrent.Add(-time.Second)

	startOfPast2 := startOfCurrent.AddDate(0, -2, 0)
	endOfPast2 := startOfPast1.Add(-time.Second)

	startOfPast3 := startOfCurrent.AddDate(0, -3, 0)
	endOfPast3 := startOfPast2.Add(-time.Second)

	startOfFuture1 := startOfCurrent.AddDate(0, 1, 0)
	endOfFuture1 := startOfFuture1.AddDate(0, 1, 0).Add(-time.Second)

	log.Println("Inserting mock ladders...")
	var past3ID, past2ID, past1ID, currentID, futureID int64

	// Past 3
	err := pool.QueryRow(ctx, `
		INSERT INTO ladders (name, type, start_time, end_time, initial_balance, is_active)
		VALUES ($1, 'monthly', $2, $3, 10000.0, false)
		RETURNING id`,
		fmt.Sprintf("Radiant %s Ladder '%02d", startOfPast3.Month().String(), startOfPast3.Year()%100),
		startOfPast3, endOfPast3).Scan(&past3ID)
	if err != nil {
		return fmt.Errorf("failed to insert past3 ladder: %w", err)
	}

	// Past 2
	err = pool.QueryRow(ctx, `
		INSERT INTO ladders (name, type, start_time, end_time, initial_balance, is_active)
		VALUES ($1, 'monthly', $2, $3, 10000.0, false)
		RETURNING id`,
		fmt.Sprintf("Breezy %s Ladder '%02d", startOfPast2.Month().String(), startOfPast2.Year()%100),
		startOfPast2, endOfPast2).Scan(&past2ID)
	if err != nil {
		return fmt.Errorf("failed to insert past2 ladder: %w", err)
	}

	// Past 1
	err = pool.QueryRow(ctx, `
		INSERT INTO ladders (name, type, start_time, end_time, initial_balance, is_active)
		VALUES ($1, 'monthly', $2, $3, 10000.0, false)
		RETURNING id`,
		fmt.Sprintf("Emerald %s Ladder '%02d", startOfPast1.Month().String(), startOfPast1.Year()%100),
		startOfPast1, endOfPast1).Scan(&past1ID)
	if err != nil {
		return fmt.Errorf("failed to insert past1 ladder: %w", err)
	}

	// Current (Active)
	err = pool.QueryRow(ctx, `
		INSERT INTO ladders (name, type, start_time, end_time, initial_balance, is_active)
		VALUES ($1, 'monthly', $2, $3, 10000.0, true)
		RETURNING id`,
		fmt.Sprintf("Active %s Ladder '%02d", startOfCurrent.Month().String(), startOfCurrent.Year()%100),
		startOfCurrent, endOfCurrent).Scan(&currentID)
	if err != nil {
		return fmt.Errorf("failed to insert active ladder: %w", err)
	}

	// Future (Upcoming)
	err = pool.QueryRow(ctx, `
		INSERT INTO ladders (name, type, start_time, end_time, initial_balance, is_active)
		VALUES ($1, 'monthly', $2, $3, 10000.0, false)
		RETURNING id`,
		fmt.Sprintf("Upcoming %s Ladder '%02d", startOfFuture1.Month().String(), startOfFuture1.Year()%100),
		startOfFuture1, endOfFuture1).Scan(&futureID)
	if err != nil {
		return fmt.Errorf("failed to insert future ladder: %w", err)
	}

	// Insert Tickers for all ladders
	ladders := []int64{past3ID, past2ID, past1ID, currentID, futureID}
	for _, lID := range ladders {
		if err := insertTickers(ctx, pool, lID); err != nil {
			return fmt.Errorf("failed to insert tickers for ladder %d: %w", lID, err)
		}
	}

	log.Println("Seeding historical participants...")
	if err := insertPastParticipants(ctx, pool, past3ID, userIDs); err != nil {
		return fmt.Errorf("failed to seed past3 participants: %w", err)
	}
	if err := insertPastParticipants(ctx, pool, past2ID, userIDs); err != nil {
		return fmt.Errorf("failed to seed past2 participants: %w", err)
	}
	if err := insertPastParticipants(ctx, pool, past1ID, userIDs); err != nil {
		return fmt.Errorf("failed to seed past1 participants: %w", err)
	}

	log.Println("Seeding active participants and portfolios...")
	if err := insertActiveParticipants(ctx, pool, currentID, userIDs); err != nil {
		return fmt.Errorf("failed to seed active participants: %w", err)
	}

	return nil
}

func insertTickers(ctx context.Context, pool *pgxpool.Pool, ladderID int64) error {
	tickers := []struct {
		symbol string
		source string
	}{
		{"AAPL", "Finnhub"},
		{"AMZN", "Finnhub"},
		{"MSFT", "Finnhub"},
		{"TSLA", "Finnhub"},
		{"bitcoin", "CoinGecko"},
		{"ethereum", "CoinGecko"},
	}
	for _, t := range tickers {
		_, err := pool.Exec(ctx, `
			INSERT INTO ladder_tickers (ladder_id, stock_symbol, source)
			VALUES ($1, $2, $3)`, ladderID, t.symbol, t.source)
		if err != nil {
			return err
		}
	}

	return nil
}

func insertPastParticipants(ctx context.Context, pool *pgxpool.Pool, ladderID int64, userIDs []int64) error {
	balances := []float64{15430.50, 14210.00, 12850.10, 11020.00, 10500.00, 9800.00, 9500.00, 8900.00, 8200.00, 7500.00}
	for i, userID := range userIDs {
		rank := i + 1
		balance := balances[i%len(balances)]
		_, err := pool.Exec(ctx, `
			INSERT INTO ladder_participants (ladder_id, user_id, balance, final_balance, final_rank)
			VALUES ($1, $2, 10000.0, $3, $4)`, ladderID, userID, balance, rank)
		if err != nil {
			return err
		}
	}

	return nil
}

func insertActiveParticipants(ctx context.Context, pool *pgxpool.Pool, ladderID int64, userIDs []int64) error {
	for _, userID := range userIDs {
		_, err := pool.Exec(ctx, `
			INSERT INTO ladder_participants (ladder_id, user_id, balance)
			VALUES ($1, $2, 10000.0)`, ladderID, userID)
		if err != nil {
			return err
		}
	}

	holdings := []struct {
		userIndex int
		symbol    string
		qty       float64
		avgPr     float64
	}{
		{1, "AAPL", 30.0, 180.00},
		{1, "bitcoin", 0.1, 60000.00},
		{2, "TSLA", 25.0, 160.00},
		{4, "MSFT", 10.0, 400.00},
		{4, "ethereum", 1.0, 3000.00},
	}

	for _, h := range holdings {
		userID := userIDs[h.userIndex]
		cost := h.qty * h.avgPr
		_, err := pool.Exec(ctx, `
			UPDATE ladder_participants
			SET balance = balance - $1
			WHERE ladder_id = $2 AND user_id = $3`, cost, ladderID, userID)
		if err != nil {
			return err
		}

		_, err = pool.Exec(ctx, `
			INSERT INTO ladder_portfolio_items (ladder_id, user_id, stock_symbol, quantity, average_price)
			VALUES ($1, $2, $3, $4, $5)`, ladderID, userID, h.symbol, h.qty, h.avgPr)
		if err != nil {
			return err
		}
	}

	return nil
}
