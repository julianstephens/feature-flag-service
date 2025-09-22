package storage

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdStore struct {
	Client    *clientv3.Client
	KeyPrefix string
}

func NewEtcdStore(endpoints []string, keyPrefix string) (*EtcdStore, error) {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 30 * time.Second,
	}
	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	return &EtcdStore{
		Client:    client,
		KeyPrefix: keyPrefix,
	}, nil
}

func (e *EtcdStore) Connect() error {
	// The etcd client connects on creation, so we can just return nil here.
	return nil
}

func (e *EtcdStore) Close() error {
	return e.Client.Close()
}

func (e *EtcdStore) List(ctx context.Context, key string, opts ...clientv3.OpOption) (map[string]string, error) {
	resp, err := e.Client.Get(ctx, key, append([]clientv3.OpOption{clientv3.WithPrefix()}, opts...)...)
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, ErrKeyNotFound
	}

	result := make(map[string]string)
	for _, kv := range resp.Kvs {
		result[string(kv.Key)] = string(kv.Value)
	}
	return result, nil
}

func (e *EtcdStore) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (string, error) {
	resp, err := e.Client.Get(ctx, key, opts...)
	if err != nil {
		return "", err
	}

	if len(resp.Kvs) == 0 {
		return "", ErrKeyNotFound
	}

	return string(resp.Kvs[0].Value), nil
}

func (e *EtcdStore) Put(ctx context.Context, key, value string, opts ...clientv3.OpOption) (string, error) {
	_, err := e.Client.Put(ctx, key, value, opts...)
	return "", err
}

func (e *EtcdStore) Post(ctx context.Context, key, value string, opts ...clientv3.OpOption) (string, error) {
	return "", ErrNotImplemented
}

func (e *EtcdStore) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) error {
	resp, err := e.Client.Delete(ctx, key, opts...)
	if err != nil {
		return err
	}
	if resp.Deleted == 0 {
		return ErrKeyNotFound
	}
	return nil
}
