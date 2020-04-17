package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

type redisStorage struct {
	client    *redis.Client
	namespace string
}

// Redis initializes a redis cache implementation wrapper
func Redis(client *redis.Client, namespace string) (Storage, error) {
	// Check arguments
	if client == nil {
		return nil, fmt.Errorf("redis client must not be nil")
	}
	if namespace == "" {
		return nil, fmt.Errorf("cache namespace must not be blank")
	}

	// Return wrapper
	return &redisStorage{
		client:    client,
		namespace: namespace,
	}, nil
}

// -----------------------------------------------------------------------------

func (s *redisStorage) Get(ctx context.Context, key string) ([]byte, error) {
	value, err := s.client.WithContext(ctx).Get(s.key(key)).Bytes()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve '%q': %w", key, err)
	}
	return value, nil
}

func (s *redisStorage) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	err := s.client.WithContext(ctx).Set(s.key(key), value, expiration).Err()
	if err != nil {
		return fmt.Errorf("unable to set '%q' value: %w", key, err)
	}
	return nil
}

func (s *redisStorage) Remove(ctx context.Context, key string) error {
	err := s.client.WithContext(ctx).Del(s.key(key)).Err()
	if err != nil {
		return fmt.Errorf("unable to remove '%q' value: %w", key, err)
	}
	return nil
}

func (s *redisStorage) key(name string) string {
	return fmt.Sprintf("%s:%s", s.namespace, name)
}
