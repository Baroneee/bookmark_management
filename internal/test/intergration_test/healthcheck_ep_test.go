package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Baroneee/EBVN_Assignment_1/internal/api"
	"github.com/Baroneee/EBVN_Assignment_1/internal/service"
	redisPkg "github.com/Baroneee/EBVN_Assignment_1/pkg/redis"
	"github.com/go-playground/assert/v2"
)

// TestHealthCheckEndpoint verifies the health check endpoint response.
func TestHealthCheckEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		setupTestRequest func() (*httptest.ResponseRecorder, *http.Request)
		expectedCode     int
		expectedBody     service.HealthResponse
	}{
		{
			name: "Case 1: Handler successfully returns health check data from service",
			setupTestRequest: func() (*httptest.ResponseRecorder, *http.Request) {
				req := httptest.NewRequest("GET", "/health", nil)
				rec := httptest.NewRecorder()
				return rec, req
			},
			expectedCode: http.StatusOK,
			expectedBody: service.HealthResponse{
				Message:     "OK",
				ServiceName: "bookmark-service",
				InstanceID:  "test-instance-id",
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testAPIEngine := api.NewEngine(&api.Config{
				ServiceName: "bookmark-service",
				InstanceID:  "test-instance-id",
			},
				redisPkg.InitMockRedis(t),
			)
			rec, req := tc.setupTestRequest()
			testAPIEngine.ServerHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			var actualBody service.HealthResponse
			err := json.Unmarshal(rec.Body.Bytes(), &actualBody)
			if err != nil {
				t.Fatalf("Failed to parse response body: %v", err)
			}

			assert.Equal(t, tc.expectedBody.Message, actualBody.Message)
			assert.Equal(t, tc.expectedBody.ServiceName, actualBody.ServiceName)
			assert.Equal(t, tc.expectedBody.InstanceID, actualBody.InstanceID)
		})
	}
}
