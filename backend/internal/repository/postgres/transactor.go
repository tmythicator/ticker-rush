package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tmythicator/ticker-rush/server/internal/service"
)

// PgxTransactor implements the Transactor interface using pgx.
type PgxTransactor struct {
	pool *pgxpool.Pool
}

// NewPgxTransactor creates a new instance of PgxTransactor.
func NewPgxTransactor(pool *pgxpool.Pool) *PgxTransactor {
	return &PgxTransactor{pool: pool}
}

// Begin starts a new transaction.
func (t *PgxTransactor) Begin(ctx context.Context) (service.Transaction, error) {
	return t.pool.Begin(ctx)
}
