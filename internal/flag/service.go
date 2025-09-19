package flag

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	ffpb "github.com/julianstephens/feature-flag-service/gen/go/grpc/v1/featureflag.v1"
	"github.com/julianstephens/feature-flag-service/internal/config"
	"github.com/julianstephens/feature-flag-service/internal/storage"
	"github.com/julianstephens/feature-flag-service/internal/utils"
)

var ErrFlagNotFound = errors.New("flag not found")

type Flag struct {
	ID          string
	Name        string
	Description string
	Enabled     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Service interface {
	CreateFlag(ctx context.Context, name, description string, enabled bool) (*Flag, error)
	UpdateFlag(ctx context.Context, id, name, description string, enabled bool) (*Flag, error)
	GetFlag(ctx context.Context, id string) (*Flag, error)
	DeleteFlag(ctx context.Context, id string) error
	ListFlags(ctx context.Context) ([]*Flag, error)
}

type EtcdService struct {
	conf  *config.Config
	store *storage.EtcdClient
}

func NewService(conf *config.Config, etcdClient *storage.EtcdClient) Service {
	return &EtcdService{
		conf:  conf,
		store: etcdClient,
	}
}

func (s *EtcdService) ListFlags(ctx context.Context) ([]*Flag, error) {
	res, err := s.store.Client.Get(ctx, s.store.KeyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var flags []*Flag
	for _, kv := range res.Kvs {
		var flag Flag
		if err := json.Unmarshal(kv.Value, &flag); err != nil {
			log.Printf("Failed to unmarshal flag data: %v", err)
			continue
		}
		flags = append(flags, &flag)
	}
	return flags, nil
}

func (s *EtcdService) GetFlag(ctx context.Context, id string) (*Flag, error) {
	resp, err := s.store.Client.Get(ctx, s.store.GetKey(id))
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, ErrFlagNotFound
	}

	var flag Flag
	if err := json.Unmarshal(resp.Kvs[0].Value, &flag); err != nil {
		return nil, err
	}
	return &flag, nil
}

func (s *EtcdService) CreateFlag(ctx context.Context, name, description string, enabled bool) (*Flag, error) {
	id := utils.GenerateID()
	now := time.Now()
	flag := &Flag{
		ID:          id,
		Name:        name,
		Description: description,
		Enabled:     enabled,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	data, err := json.Marshal(flag)
	if err != nil {
		return nil, err
	}

	_, err = s.store.Client.Put(ctx, s.store.GetKey(id), string(data))
	if err != nil {
		return nil, err
	}

	return flag, nil
}

func (s *EtcdService) UpdateFlag(ctx context.Context, id, name, description string, enabled bool) (*Flag, error) {
	flag, err := s.GetFlag(ctx, id)
	if err != nil {
		return nil, err
	}

	flag.Name = name
	flag.Description = description
	flag.Enabled = enabled
	flag.UpdatedAt = time.Now()

	data, err := json.Marshal(flag)
	if err != nil {
		return nil, err
	}

	_, err = s.store.Client.Put(ctx, s.store.GetKey(id), string(data))
	if err != nil {
		return nil, err
	}

	return flag, nil
}

func (s *EtcdService) DeleteFlag(ctx context.Context, id string) error {
	res, err := s.store.Client.Delete(ctx, s.store.GetKey(id))
	if err != nil {
		return err
	}
	if res.Deleted == 0 {
		return ErrFlagNotFound
	}
	return nil
}


func (f *Flag) ToProto() *ffpb.Flag {
	return &ffpb.Flag{
		Id:          f.ID,
		Name:        f.Name,
		Description: f.Description,
		Enabled:     f.Enabled,
		CreatedAt:   f.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   f.UpdatedAt.Format(time.RFC3339),
	}
}

func FlagFromProto(protoFlag *ffpb.Flag) (*Flag, error) {
	createdAt, err := time.Parse(time.RFC3339, protoFlag.CreatedAt)
	if err != nil {
		return nil, err
	}
	updatedAt, err := time.Parse(time.RFC3339, protoFlag.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &Flag{
		ID:          protoFlag.Id,
		Name:        protoFlag.Name,
		Description: protoFlag.Description,
		Enabled:     protoFlag.Enabled,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}