package repository

import (
	"context"
	"testing"

	redisPkg "github.com/Baroneee/EBVN_Assignment_1/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// TestUrlStorage_StoreURL verifies that StoreURL persists a URL under the
// requested key and returns the expected error result.
func TestUrlStorage_StoreURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		setupMock   func(ctx context.Context) *redis.Client
		expectedErr error
		verifyFunc  func(ctx context.Context, r *redis.Client)
	}{
		{
			name: "normal case",
			setupMock: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				res, err := r.Get(ctx, "abcxyz").Result()
				assert.NoError(t, err)
				assert.Equal(t, "google.com", res)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mock := tc.setupMock(ctx)
			testRepo := NewURLStorage(mock)

			err := testRepo.StoreURL(ctx, "abcxyz", "google.com")
			assert.ErrorIs(t, err, tc.expectedErr)
			if tc.verifyFunc != nil {
				tc.verifyFunc(ctx, mock)
			}
		})
	}
}

// TestURLStorage_StoreURLIfAbsent verifies that an existing key is never
// overwritten by a subsequent write.
func TestURLStorage_StoreURLIfAbsent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		setupMock   func(ctx context.Context) *redis.Client
		expectedErr error
		verifyFunc  func(ctx context.Context, r *redis.Client)
	}{
		{
			name: "existing key is not overwritten",
			setupMock: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				res, err := r.Get(ctx, "abcxyz").Result()
				assert.NoError(t, err)
				assert.Equal(t, "google.com", res)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mock := tc.setupMock(ctx)
			testRepo := NewURLStorage(mock)

			created, err := testRepo.StoreURLIfAbsent(ctx, "abcxyz", "google.com")
			assert.NoError(t, err)
			assert.True(t, created)

			created, err = testRepo.StoreURLIfAbsent(ctx, "abcxyz", "example.com")
			assert.ErrorIs(t, err, tc.expectedErr)
			assert.False(t, created)

			if tc.verifyFunc != nil {
				tc.verifyFunc(ctx, mock)
			}
		})
	}
}
