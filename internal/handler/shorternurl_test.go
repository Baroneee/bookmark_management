package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	svc_mocks "github.com/Baroneee/EBVN_Assignment_1/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestShortenURL_CreateShorternLink verifies the shortened-link handler response.
func TestShortenURL_CreateShorternLink(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		body             string
		setupMockService func(*svc_mocks.ShortenURL)
		expectedCode     int
		expectedBody     map[string]string
	}{
		{
			name: "success",
			body: `{"url":"https://example.com"}`,
			setupMockService: func(mockSvc *svc_mocks.ShortenURL) {
				mockSvc.On("CreateShorternLink", mock.Anything, "https://example.com").
					Return("abc1234", nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]string{
				"code": "abc1234",
			},
		},
		{
			name:             "invalid input",
			body:             `{`,
			setupMockService: func(mockSvc *svc_mocks.ShortenURL) {},
			expectedCode:     http.StatusBadRequest,
			expectedBody:     map[string]string{"error": "Invalid Input"},
		},
		{
			name: "service error",
			body: `{"url":"https://example.com"}`,
			setupMockService: func(mockSvc *svc_mocks.ShortenURL) {
				mockSvc.On("CreateShorternLink", mock.Anything, "https://example.com").
					Return("", assert.AnError)
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]string{"error": "Internal Server Error"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shortern",
				bytes.NewBufferString(tc.body))
			ctx.Request.Header.Set("Content-Type", "application/json")

			mockSvc := &svc_mocks.ShortenURL{}
			tc.setupMockService(mockSvc)
			shortenHandler := NewShortenURL(mockSvc)
			shortenHandler.CreateShorternLink(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			var actualBody map[string]string
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actualBody))
			assert.Equal(t, tc.expectedBody, actualBody)
			mockSvc.AssertExpectations(t)
		})
	}
}
