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
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Service interface {
	CreateFlag(ctx context.Context, name, description string, enabled bool) (*Flag, error)
	UpdateFlag(ctx context.Context, id, name, description string, enabled bool) (*Flag, error)
	GetFlag(ctx context.Context, id string) (*Flag, error)
	DeleteFlag(ctx context.Context, id string) error
	ListFlags(ctx context.Context) ([]*Flag, error)
}

type FlagService struct {
	conf   *config.Config
	store  storage.Store[clientv3.OpOption]
	prefix string
}

func NewService(conf *config.Config, etcdClient *storage.EtcdStore) Service {
	return &FlagService{
		conf:   conf,
		store:  etcdClient,
		prefix: conf.FlagServicePrefix,
	}
}

func (s *FlagService) GetKey(key string) string {
	return s.prefix + key
}

func (s *FlagService) ListFlags(ctx context.Context) ([]*Flag, error) {
	res, err := s.store.List(ctx, "")
	if err != nil {
		return nil, err
	}

	var flags []*Flag
	for _, v := range res {
		var flag Flag
		if err := json.Unmarshal([]byte(v), &flag); err != nil {
			log.Printf("error unmarshaling flag: %v", err)
			continue
		}
		flags = append(flags, &flag)
	}
	return flags, nil
}

func (s *FlagService) GetFlag(ctx context.Context, id string) (*Flag, error) {
	resp, err := s.store.Get(ctx, s.GetKey(id))
	if err != nil {
		if errors.Is(err, storage.ErrKeyNotFound) {
			return nil, ErrFlagNotFound
		}
		return nil, err
	}

	var flag Flag
	if err := json.Unmarshal([]byte(resp), &flag); err != nil {
		return nil, err
	}
	return &flag, nil
}

func (s *FlagService) CreateFlag(ctx context.Context, name, description string, enabled bool) (*Flag, error) {
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

	_, err = s.store.Put(ctx, s.GetKey(id), string(data))
	if err != nil {
		return nil, err
	}

	return flag, nil
}

func (s *FlagService) UpdateFlag(ctx context.Context, id, name, description string, enabled bool) (*Flag, error) {
	flag, err := s.GetFlag(ctx, s.GetKey(id))
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

	_, err = s.store.Put(ctx, s.GetKey(id), string(data))
	if err != nil {
		return nil, err
	}

	return flag, nil
}

func (s *FlagService) DeleteFlag(ctx context.Context, id string) error {
	return s.store.Delete(ctx, s.GetKey(id))
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

func ParseFlag(data []byte) (*Flag, error) {
	var flag Flag
	if err := json.Unmarshal(data, &flag); err != nil {
		return nil, err
	}
	return &flag, nil
}
