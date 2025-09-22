package storage

import (
	"context"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5"

	"github.com/julianstephens/feature-flag-service/internal/config"
)

var (
	dbConn *pgx.Conn
	once   sync.Once
)

type PostgresStore struct {
	TableName string
	Columns   []string
	IdxKey    string
}

type PostgresOption struct {
	TableName string
	Columns   []string
	IdxKey    string
}

func NewPostgresStore(opts PostgresOption) *PostgresStore {
	return &PostgresStore{
		TableName: opts.TableName,
		Columns:   opts.Columns,
		IdxKey:    opts.IdxKey,
	}
}

func (s *PostgresStore) Connect() error {
	var err error
	once.Do(func() {
		conf := config.LoadConfig()
		conn, e := pgx.Connect(context.Background(), conf.PostgresURL)
		if e != nil {
			err = e
		}
		dbConn = conn
	})
	return err
}

func (s *PostgresStore) Close() error {
	if dbConn != nil {
		return dbConn.Close(context.Background())
	}
	return nil
}

func (s *PostgresStore) List(ctx context.Context, prefix string, opts ...any) (map[string]string, error) {
	if dbConn == nil {
		if err := s.Connect(); err != nil {
			return nil, err
		}
	}

	query := "SELECT " + strings.Join(s.Columns, ",") + " FROM " + s.TableName
	args := []interface{}{}

	if prefix != "" {
		query += " WHERE " + s.IdxKey + " LIKE $1"
		args = append(args, prefix+"%")
	}

	rows, err := dbConn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		result[key] = value
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PostgresStore) Get(ctx context.Context, key string, opts ...any) (string, error) {
	if dbConn == nil {
		if err := s.Connect(); err != nil {
			return "", err
		}
	}

	query := "SELECT " + strings.Join(s.Columns, ",") + " FROM " + s.TableName + " WHERE " + s.IdxKey + "=$1"
	row := dbConn.QueryRow(ctx, query, key)

	var result string
	if err := row.Scan(&result); err != nil {
		if err == pgx.ErrNoRows {
			return "", ErrKeyNotFound
		}
		return "", err
	}
	return result, nil
}

func (s *PostgresStore) Post(ctx context.Context, key, value string, opts ...any) (string, error) {
	return "", nil
}

func (s *PostgresStore) Put(ctx context.Context, key, value string, opts ...any) (string, error) {
	return "", nil
}

func (s *PostgresStore) Delete(ctx context.Context, key string, opts ...any) error {
	return nil
}
