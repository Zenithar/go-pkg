package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/allegro/bigcache"
)

type bcStorage struct {
	store *bigcache.BigCache
}

// BigCache initializes a bigcache implementation wrapper
func BigCache(cfg bigcache.Config) (Storage, error) {
	// Initialize bigcache backend
	store, err := bigcache.NewBigCache(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize cache: %w", err)
	}

	// Return wrapper
	return &bcStorage{
		store: store,
	}, nil
}

// -----------------------------------------------------------------------------

func (s *bcStorage) Get(_ context.Context, key string) ([]byte, error) {
	value, err := s.store.Get(key)
	if err != nil {
		if err == bigcache.ErrEntryNotFound {
			return nil, ErrCacheMiss
		}
		return nil, fmt.Errorf("unable to retrieve '%q': %w", key, err)
	}
	return value, nil
}

func (s *bcStorage) Set(_ context.Context, key string, value []byte, _ time.Duration) error {
	err := s.store.Set(key, value)
	if err != nil {
		return fmt.Errorf("unable to set '%q' value: %w", key, err)
	}
	return nil
}

func (s *bcStorage) Remove(ctx context.Context, key string) error {
	err := s.store.Delete(key)
	if err != nil {
		return fmt.Errorf("unable to remove '%q' value: %w", key, err)
	}
	return nil
}
