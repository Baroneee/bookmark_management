package repository

import (
	"context"
	"testing"
	"time"

	redisPkg "github.com/Baroneee/EBVN_Assignment_1/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// TestURLStorage verifies URL persistence, lookup, expiry, and Redis errors.
func TestURLStorage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		operation string
		setup     func(context.Context) *redis.Client
	}{
		{
			name:      "store URL",
			operation: "store",
			setup: func(context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
		},
		{
			name:      "store URL error",
			operation: "store-error",
			setup: func(context.Context) *redis.Client {
				return redis.NewClient(&redis.Options{Addr: "127.0.0.1:1"})
			},
		},
		{
			name:      "store if absent",
			operation: "store-if-absent",
			setup: func(context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
		},
		{
			name:      "store if absent default expiry",
			operation: "store-if-absent-default",
			setup: func(context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
		},
		{
			name:      "store if absent error",
			operation: "store-if-absent-error",
			setup: func(context.Context) *redis.Client {
				return redis.NewClient(&redis.Options{Addr: "127.0.0.1:1"})
			},
		},
		{
			name:      "get URL",
			operation: "get",
			setup: func(ctx context.Context) *redis.Client {
				rdb := redisPkg.InitMockRedis(t)
				assert.NoError(t, rdb.Set(ctx, "abc1234", "https://example.com", time.Hour).Err())
				return rdb
			},
		},
		{
			name:      "get missing URL",
			operation: "get-missing",
			setup: func(context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
		},
		{
			name:      "get URL error",
			operation: "get-error",
			setup: func(context.Context) *redis.Client {
				return redis.NewClient(&redis.Options{Addr: "127.0.0.1:1"})
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			rdb := tc.setup(ctx)
			testRepo := NewURLStorage(rdb)

			switch tc.operation {
			case "store":
				assert.NoError(t, testRepo.StoreURL(ctx, "abcxyz", "google.com"))
				result, err := testRepo.GetURL(ctx, "abcxyz")
				assert.NoError(t, err)
				assert.Equal(t, "google.com", result)
			case "store-error":
				assert.Error(t, testRepo.StoreURL(ctx, "abcxyz", "google.com"))
			case "store-if-absent":
				created, err := testRepo.StoreURLIfAbsent(ctx, "abcxyz", "google.com", time.Hour)
				assert.NoError(t, err)
				assert.True(t, created)
				created, err = testRepo.StoreURLIfAbsent(ctx, "abcxyz", "example.com", time.Hour)
				assert.NoError(t, err)
				assert.False(t, created)
				result, err := testRepo.GetURL(ctx, "abcxyz")
				assert.NoError(t, err)
				assert.Equal(t, "google.com", result)
			case "store-if-absent-default":
				created, err := testRepo.StoreURLIfAbsent(ctx, "abc1234", "https://example.com", 0)
				assert.NoError(t, err)
				assert.True(t, created)
				ttl, err := rdb.TTL(ctx, "abc1234").Result()
				assert.NoError(t, err)
				assert.Greater(t, ttl, time.Duration(0))
			case "store-if-absent-error":
				created, err := testRepo.StoreURLIfAbsent(ctx, "abc1234", "https://example.com", time.Hour)
				assert.Error(t, err)
				assert.False(t, created)
			case "get":
				result, err := testRepo.GetURL(ctx, "abc1234")
				assert.NoError(t, err)
				assert.Equal(t, "https://example.com", result)
			case "get-missing":
				result, err := testRepo.GetURL(ctx, "missing")
				assert.ErrorIs(t, err, ErrUrlDoesNotExist)
				assert.Empty(t, result)
			case "get-error":
				result, err := testRepo.GetURL(ctx, "abc1234")
				assert.Error(t, err)
				assert.Empty(t, result)
			}
		})
	}
}
