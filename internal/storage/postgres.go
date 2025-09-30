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
func (p *PostgresStore) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.db.Query(ctx, sql, args...)
}

// Exec helper
func (p *PostgresStore) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := p.db.Exec(ctx, sql, args...)
	return err
}

// ListAll retrieves all rows from the specified table
func (p *PostgresStore) ListAll(ctx context.Context, table string, dest any) error {
	sql := fmt.Sprintf("SELECT * FROM %s", table)
	if err := pgxscan.Select(ctx, p.db, dest, sql); err != nil {
		return err
	}
	return nil
}

func (p *PostgresStore) Get(ctx context.Context, table string, dest any, whereClause string, args ...any) error {
	sql := fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT 1", table, whereClause)
	if err := pgxscan.Get(ctx, p.db, dest, sql, args...); err != nil {
		return err
	}
	return nil
}

func (p *PostgresStore) Post(ctx context.Context, table string, data map[string]any) error {
	columns := ""
	values := ""
	args := []any{}
	i := 1

	for k, v := range data {
		if columns != "" {
			columns += ", "
			values += ", "
		}
		columns += k
		values += fmt.Sprintf("$%d", i)
		args = append(args, v)
		i++
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, columns, values)
	_, err := p.db.Exec(ctx, sql, args...)
	return err
}

func (p *PostgresStore) Put(ctx context.Context, table string, data map[string]any, whereClause string, whereArgs ...any) error {
	setClause := ""
	args := []any{}
	i := 1

	for k, v := range data {
		if setClause != "" {
			setClause += ", "
		}
		setClause += fmt.Sprintf("%s=$%d", k, i)
		args = append(args, v)
		i++
	}

	args = append(args, whereArgs...)
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, setClause, whereClause)
	_, err := p.db.Exec(ctx, sql, args...)
	return err
}
