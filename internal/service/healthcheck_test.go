package service

import (
	"testing"

	redisPkg "github.com/Baroneee/EBVN_Assignment_1/pkg/redis"
	"github.com/go-playground/assert/v2"
	"github.com/redis/go-redis/v9"
)

// TestHealthCheck verifies health status creation and instance ID behavior.
func TestHealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		serviceName    string
		instanceID     string
		redisClient    *redis.Client
		expectedMsg    string
		expectedName   string
		expectedErr    bool
		checkEmptyID   bool
		expectedCustom string
	}{
		{
			name:         "Case 1: InstanceID empty -> Should auto-generate UUID",
			serviceName:  "bookmark-service",
			instanceID:   "",
			redisClient:  redisPkg.InitMockRedis(t),
			expectedMsg:  "OK",
			expectedName: "bookmark-service",
			checkEmptyID: true,
		},
		{
			name:           "Case 2: InstanceID provided -> Should use provided ID",
			serviceName:    "custom-service",
			instanceID:     "my-fixed-instance-id-123",
			redisClient:    redisPkg.InitMockRedis(t),
			expectedMsg:    "OK",
			expectedName:   "custom-service",
			expectedCustom: "my-fixed-instance-id-123",
		},
		{
			name:         "Case 3: Redis client unavailable",
			serviceName:  "bookmark-service",
			instanceID:   "instance-id",
			expectedName: "bookmark-service",
			expectedErr:  true,
		},
		{
			name:         "Case 4: Redis check fails",
			serviceName:  "bookmark-service",
			instanceID:   "instance-id",
			redisClient:  redis.NewClient(&redis.Options{Addr: "127.0.0.1:1"}),
			expectedName: "bookmark-service",
			expectedErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			svc := NewHealthCheck(tc.serviceName, tc.instanceID, tc.redisClient)
			result, err := svc.Check()

			if tc.expectedErr {
				assert.NotEqual(t, nil, err)
			} else {
				assert.Equal(t, nil, err)
			}
			assert.Equal(t, tc.expectedMsg, result.Message)
			assert.Equal(t, tc.expectedName, result.ServiceName)
			if tc.checkEmptyID {
				assert.NotEqual(t, "", result.InstanceID)
			}
			if tc.expectedCustom != "" {
				assert.Equal(t, tc.expectedCustom, result.InstanceID)
			}
		})
	}
}
