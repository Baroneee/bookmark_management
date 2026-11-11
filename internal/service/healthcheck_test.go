package service

import (
	"testing"

	redisPkg "github.com/Baroneee/EBVN_Assignment_1/pkg/redis"
	"github.com/go-playground/assert/v2"
)

// TestHealthCheck verifies health status creation and instance ID behavior.
func TestHealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		serviceName    string
		instanceID     string
		expectedMsg    string
		expectedName   string
		checkEmptyID   bool
		expectedCustom string
	}{
		{
			name:         "Case 1: InstanceID empty -> Should auto-generate UUID",
			serviceName:  "bookmark-service",
			instanceID:   "",
			expectedMsg:  "OK",
			expectedName: "bookmark-service",
			checkEmptyID: true,
		},
		{
			name:           "Case 2: InstanceID provided -> Should use provided ID",
			serviceName:    "custom-service",
			instanceID:     "my-fixed-instance-id-123",
			expectedMsg:    "OK",
			expectedName:   "custom-service",
			expectedCustom: "my-fixed-instance-id-123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewHealthCheck(tc.serviceName, tc.instanceID, redisPkg.InitMockRedis(t))
			result, err := svc.Check()

			assert.Equal(t, nil, err)
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
