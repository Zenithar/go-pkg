package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"golang.org/x/xerrors"
)

type redisStorage struct {
	client    *redis.Client
	namespace string
}

// Redis initializes a redis cache implementation wrapper
func Redis(client *redis.Client, namespace string) (Storage, error) {
	// Return wrapper
	return &redisStorage{
		client:    client,
		namespace: namespace,
	}, nil
}

// -----------------------------------------------------------------------------

func (s *redisStorage) Get(_ context.Context, key string) ([]byte, error) {
	value, err := s.client.Get(s.key(key)).Bytes()
	if err != nil {
		return nil, xerrors.Errorf("bigcache: unable to retrieve '%q': %w", key, err)
	}
	return value, nil
}

func (s *redisStorage) Set(_ context.Context, key string, value []byte, expiration time.Duration) error {
	err := s.client.Set(s.key(key), value, expiration).Err()
	if err != nil {
		return xerrors.Errorf("bigcache: unable to set '%q' value: %w", key, err)
	}
	return nil
}

func (s *redisStorage) Remove(_ context.Context, key string) error {
	err := s.client.Del(s.key(key)).Err()
	if err != nil {
		return xerrors.Errorf("bigcache: unable to remove '%q' value: %w", key, err)
	}
	return nil
}

func (s *redisStorage) key(name string) string {
	return fmt.Sprintf("%s:%s", s.namespace, name)
}
