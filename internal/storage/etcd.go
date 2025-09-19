package storage

import (
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)


type EtcdClient struct {
	Client *clientv3.Client
	KeyPrefix string	
}

func NewEtcdClient(endpoints []string, keyPrefix string) (*EtcdClient, error) {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	}
	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	return &EtcdClient{
		Client: client,
		KeyPrefix: keyPrefix,
	}, nil
}

func (e *EtcdClient) GetKey(key string) string {
	return e.KeyPrefix + key
}