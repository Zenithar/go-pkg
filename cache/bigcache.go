package cache

import (
	"context"
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
		return nil, err
	}

	// Return wrapper
	return &bcStorage{
		store: store,
	}, nil
}

// -----------------------------------------------------------------------------

func (s *bcStorage) Get(_ context.Context, key string) ([]byte, error) {
	return s.store.Get(key)
}

func (s *bcStorage) Set(_ context.Context, key string, value []byte, _ time.Duration) error {
	return s.store.Set(key, value)
}

func (s *bcStorage) Remove(ctx context.Context, key string) error {
	return s.store.Delete(key)
}
