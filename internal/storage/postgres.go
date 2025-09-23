package storage

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/julianstephens/feature-flag-service/internal/config"
)

type PostgresStore struct {
    db *pgxpool.Pool
}

func NewPostgresStore(conf *config.Config) (*PostgresStore, error) {
    pool, err := pgxpool.New(context.Background(), conf.PostgresURL)
    if err != nil {
        return nil, err
    }

	// Test connection
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

    return &PostgresStore{db: pool}, nil
}

func (p *PostgresStore) Close() error {
    if p.db != nil {
        p.db.Close()
    }
    return nil
}

// Query helper
func (p *PostgresStore) Query(ctx context.Context,sql string, args ...any) (pgx.Rows, error) {
	return p.db.Query(ctx, sql, args...)
}

// Exec helper
func (p *PostgresStore) Exec(ctx context.Context, sql string, args ...any) error {
    _, err := p.db.Exec(ctx, sql, args...)
    return err
}

// ListAll retrieves all rows from the specified table
func (p *PostgresStore) ListAll(ctx context.Context, table string, dest []any) error {
    sql := fmt.Sprintf("SELECT * FROM %s", table)
	if err := pgxscan.Select(ctx, p.db, dest, sql); err != nil {
		return err
	}
	return nil
}

func (p *PostgresStore) Get(ctx context.Context, table string, dest any, where string, args ...any) error {
	sql := fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT 1", table, where)
	if err := pgxscan.Get(ctx, p.db, dest, sql, args...); err != nil {
		return err
	}
	return nil
}
