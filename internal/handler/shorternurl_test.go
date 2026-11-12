package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Baroneee/EBVN_Assignment_1/internal/service"
	svc_mocks "github.com/Baroneee/EBVN_Assignment_1/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestShortenURL verifies shortened-link creation and redirect responses.
func TestShortenURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		operation        string
		body             string
		code             string
		setupMockService func(*svc_mocks.ShortenURL)
		expectedCode     int
		expectedBody     map[string]string
		expectedURL      string
	}{
		{
			name:      "create success",
			operation: "create",
			body:      `{"url":"https://example.com","expireTime":7200}`,
			setupMockService: func(mockSvc *svc_mocks.ShortenURL) {
				mockSvc.On("CreateShorternLink", mock.Anything, "https://example.com", 2*time.Hour).
					Return("abc1234", nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]string{"code": "abc1234"},
		},
		{
			name:             "create invalid input",
			operation:        "create",
			body:             `{`,
			setupMockService: func(*svc_mocks.ShortenURL) {},
			expectedCode:     http.StatusBadRequest,
			expectedBody:     map[string]string{"error": "Invalid Input"},
		},
		{
			name:             "create invalid URL",
			operation:        "create",
			body:             `{"url":"string"}`,
			setupMockService: func(*svc_mocks.ShortenURL) {},
			expectedCode:     http.StatusBadRequest,
			expectedBody:     map[string]string{"error": "Invalid Input"},
		},
		{
			name:      "create service error",
			operation: "create",
			body:      `{"url":"https://example.com"}`,
			setupMockService: func(mockSvc *svc_mocks.ShortenURL) {
				mockSvc.On("CreateShorternLink", mock.Anything, "https://example.com", time.Duration(0)).
					Return("", assert.AnError)
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]string{"error": "Internal Server Error"},
		},
		{
			name:             "create negative expiry",
			operation:        "create",
			body:             `{"url":"https://example.com","expireTime":-1}`,
			setupMockService: func(*svc_mocks.ShortenURL) {},
			expectedCode:     http.StatusBadRequest,
			expectedBody:     map[string]string{"error": "Invalid Input"},
		},
		{
			name:             "create expiry exceeds limit",
			operation:        "create",
			body:             `{"url":"https://example.com","expireTime":604801}`,
			setupMockService: func(*svc_mocks.ShortenURL) {},
			expectedCode:     http.StatusBadRequest,
			expectedBody:     map[string]string{"error": "Invalid Input"},
		},
		{
			name:      "redirect success",
			operation: "redirect",
			code:      "abc1234",
			setupMockService: func(mockSvc *svc_mocks.ShortenURL) {
				mockSvc.On("GetUrlFromCode", mock.Anything, "abc1234").
					Return("https://example.com", nil)
			},
			expectedCode: http.StatusMovedPermanently,
			expectedURL:  "https://example.com",
		},
		{
			name:             "redirect empty code",
			operation:        "redirect",
			setupMockService: func(*svc_mocks.ShortenURL) {},
			expectedCode:     http.StatusBadRequest,
			expectedBody:     map[string]string{"error": "Invalid Input"},
		},
		{
			name:      "redirect code does not exist",
			operation: "redirect",
			code:      "missing",
			setupMockService: func(mockSvc *svc_mocks.ShortenURL) {
				mockSvc.On("GetUrlFromCode", mock.Anything, "missing").
					Return("", service.ErrUrlDoesNotExist)
			},
			expectedCode: http.StatusNotFound,
			expectedBody: map[string]string{"error": "Code Not Exist"},
		},
		{
			name:      "redirect service error",
			operation: "redirect",
			code:      "abc1234",
			setupMockService: func(mockSvc *svc_mocks.ShortenURL) {
				mockSvc.On("GetUrlFromCode", mock.Anything, "abc1234").
					Return("", errors.New("storage unavailable"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]string{"error": "Internal Server Error"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			mockSvc := &svc_mocks.ShortenURL{}
			tc.setupMockService(mockSvc)

			if tc.operation == "create" {
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shortern",
					bytes.NewBufferString(tc.body))
				ctx.Request.Header.Set("Content-Type", "application/json")
				NewShortenURL(mockSvc).CreateShorternLink(ctx)
			} else {
				ctx.Request = httptest.NewRequest(http.MethodGet,
					"/v1/links/redirect/"+tc.code, nil)
				ctx.Params = gin.Params{{Key: "code", Value: tc.code}}
				NewShortenURL(mockSvc).Redirect(ctx)
			}

			assert.Equal(t, tc.expectedCode, rec.Code)
			if tc.expectedURL != "" {
				assert.Equal(t, tc.expectedURL, rec.Header().Get("Location"))
			} else {
				var actualBody map[string]string
				assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actualBody))
				assert.Equal(t, tc.expectedBody, actualBody)
			}
			mockSvc.AssertExpectations(t)
		})
	}
}
