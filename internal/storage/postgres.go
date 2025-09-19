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
	TableName   string
	Columns     []string
	IdxKey      string
}

func NewPostgresStore(opts PostgresOption) *PostgresStore {
	return &PostgresStore{
		TableName: opts.TableName,
		Columns:   opts.Columns,
		IdxKey:    opts.IdxKey,
	}
}

func Connect() error {
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

func Close() error {
	if dbConn != nil {
		return dbConn.Close(context.Background())
	}
	return nil
}


func (s *PostgresStore) List(ctx context.Context, prefix string, opts ...any) (map[string]string, error) {
	if dbConn == nil {
		if err := Connect(); err != nil {
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

func Get(ctx context.Context, prefix string, opts ...any) (string, error) {
	return "", nil
}

func Post(ctx context.Context, prefix string, opts ...any) (string, error) {
	return "", nil
}

func Put(ctx context.Context, prefix string, opts ...any) (string, error) {
	return "", nil
}

func Delete(ctx context.Context, prefix string, opts ...any) error {
	return nil
}