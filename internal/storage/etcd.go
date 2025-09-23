package storage

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdStore struct {
    client    *clientv3.Client
    keyPrefix string
}

func NewEtcdStore(endpoints []string, prefix string) (*EtcdStore, error) {
    cli, err := clientv3.New(clientv3.Config{
        Endpoints:   endpoints,
        DialTimeout: time.Second * 10,
    })
    if err != nil {
        return nil, err
    }
    return &EtcdStore{
        client:    cli,
        keyPrefix: prefix,
    }, nil
}

// Close the etcd client connection
func (e *EtcdStore) Close() error {
    return e.client.Close()
}

// Get helper to get a single key
func (e *EtcdStore) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
    return e.client.Get(ctx, e.keyPrefix+key, opts...)
}

// Put helper to put a single key
func (e *EtcdStore) Put(ctx context.Context, key, value string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
    return e.client.Put(ctx, e.keyPrefix+key, value, opts...)
}

// Delete helper to delete a single key
func (e *EtcdStore) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
    return e.client.Delete(ctx, e.keyPrefix+key, opts...)
}

// List all keys with a given prefix
func (e *EtcdStore) ListAll(ctx context.Context, prefix string) (map[string]string, error) {
    resp, err := e.client.Get(ctx, e.keyPrefix+prefix, clientv3.WithPrefix())
    if err != nil {
        return nil, err
    }
    result := make(map[string]string)
    for _, kv := range resp.Kvs {
        result[string(kv.Key)] = string(kv.Value)
    }
    return result, nil
}

// ParseSingleGetResponse parses a single get response
func (e *EtcdStore) ParseSingleGetResponse(resp *clientv3.GetResponse) ([]byte, error) {
    if resp.Count == 0 {
        return nil, nil
    }
    return resp.Kvs[0].Value, nil
}

// ParseSingleDeleteResponse parses a single delete response
func (e *EtcdStore) ParseSingleDeleteResponse(resp *clientv3.DeleteResponse) (int64, error) {
    return resp.Deleted, nil
}

// ParseSinglePutResponse parses a single put response
func (e *EtcdStore) ParseSinglePutResponse(resp *clientv3.PutResponse) ([]byte, error) {
    if resp.PrevKv == nil {
        return nil, nil
    }
    return resp.PrevKv.Value, nil
}

// ParseListResponse parses a list response
func (e *EtcdStore) ParseListResponse(resp *clientv3.GetResponse) (map[string][]byte, error) {
       result := make(map[string][]byte)
       for _, kv := range resp.Kvs {
           result[string(kv.Key)] = kv.Value
       }
       return result, nil
}