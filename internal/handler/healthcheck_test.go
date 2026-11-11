package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Baroneee/EBVN_Assignment_1/internal/service"
	svc_mocks "github.com/Baroneee/EBVN_Assignment_1/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

// TestHealthCheck verifies the health check handler response.
func TestHealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		setupMockService func(ctx context.Context) *svc_mocks.HealthCheck
		expectedCode     int
		expectedBody     service.HealthResponse
	}{
		{
			name: "Case 1: Handler successfully returns health check data from service",
			setupMockService: func(ctx context.Context) *svc_mocks.HealthCheck {
				mockSvc := &svc_mocks.HealthCheck{}
				mockSvc.On("Check").Return(service.HealthResponse{
					Message:     "OK",
					ServiceName: "bookmark-service",
					InstanceID:  "mock-uuid-123",
				}, nil)
				return mockSvc
			},
			expectedCode: http.StatusOK,
			expectedBody: service.HealthResponse{
				Message:     "OK",
				ServiceName: "bookmark-service",
				InstanceID:  "mock-uuid-123",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(rec)
			mockSvc := tc.setupMockService(ginCtx)
			healthCheckHandler := NewHealthCheck(mockSvc)
			healthCheckHandler.Check(ginCtx)

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
