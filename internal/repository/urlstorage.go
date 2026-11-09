package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const expTime = 1 * time.Hour

// URLStorage stores and retrieves URLs by their shortened keys.
type URLStorage interface {
	StoreURL(ctx context.Context, key, url string) error
	StoreURLIfAbsent(ctx context.Context, key, url string) (bool, error)
	GetURL(ctx context.Context, key string) (string, error)
}

type urlStorage struct {
	c *redis.Client
}

// NewURLStorage creates a URL storage repository backed by Redis.
func NewURLStorage(c *redis.Client) URLStorage {
	return &urlStorage{c: c}
}

// StoreURL stores a URL under key with the repository's expiration time.
func (s *urlStorage) StoreURL(ctx context.Context, key, url string) error {
	return s.c.Set(ctx, key, url, expTime).Err()
}

// StoreURLIfAbsent stores a URL only when key does not already exist.
func (s *urlStorage) StoreURLIfAbsent(ctx context.Context, key, url string) (bool, error) {
	return s.c.SetNX(ctx, key, url, expTime).Result()
}

// GetURL retrieves the URL stored under key.
func (s *urlStorage) GetURL(ctx context.Context, key string) (string, error) {
	url, err := s.c.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return url, nil
}
