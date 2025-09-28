package storage

import (
	"context"
	"errors"
)

var (
	ErrKeyNotFound    = errors.New("key not found")
	ErrKeyExists      = errors.New("key already exists")
	ErrNotImplemented = errors.New("not implemented")
)

type Store[T any] interface {
	Connect() error
	Close() error
	List(ctx context.Context, prefix string, opts ...T) (map[string]string, error)
	Get(ctx context.Context, key string, opts ...T) (string, error)
	Post(ctx context.Context, key, value string, opts ...T) (string, error)
	Put(ctx context.Context, key, value string, opts ...T) (string, error)
	Delete(ctx context.Context, key string, opts ...T) error
}
