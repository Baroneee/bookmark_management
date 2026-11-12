package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Baroneee/EBVN_Assignment_1/internal/repository"
	repository_mocks "github.com/Baroneee/EBVN_Assignment_1/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestShortenURL verifies short URL creation, lookup, retry, and error behavior.
func TestShortenURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		operation        string
		setupMockStorage func(*repository_mocks.URLStorage)
		expectedError    error
		expectedURL      string
		verifyCode       bool
		expectedCalls    int
	}{
		{
			name:      "create success",
			operation: "create",
			setupMockStorage: func(mockStorage *repository_mocks.URLStorage) {
				mockStorage.On("StoreURLIfAbsent", mock.Anything, mock.Anything, "https://example.com", time.Hour).
					Return(true, nil).Once()
			},
			verifyCode: true, expectedCalls: 1,
		},
		{
			name:      "create retries after collision",
			operation: "create",
			setupMockStorage: func(mockStorage *repository_mocks.URLStorage) {
				mockStorage.On("StoreURLIfAbsent", mock.Anything, mock.Anything, "https://example.com", time.Hour).
					Return(false, nil).Once()
				mockStorage.On("StoreURLIfAbsent", mock.Anything, mock.Anything, "https://example.com", time.Hour).
					Return(true, nil).Once()
			},
			verifyCode: true, expectedCalls: 2,
		},
		{
			name:      "create storage error",
			operation: "create",
			setupMockStorage: func(mockStorage *repository_mocks.URLStorage) {
				mockStorage.On("StoreURLIfAbsent", mock.Anything, mock.Anything, "https://example.com", time.Hour).
					Return(false, assert.AnError).Once()
			},
			expectedError: assert.AnError, expectedCalls: 1,
		},
		{
			name:      "create exceeds maximum attempts",
			operation: "create",
			setupMockStorage: func(mockStorage *repository_mocks.URLStorage) {
				mockStorage.On("StoreURLIfAbsent", mock.Anything, mock.Anything, "https://example.com", time.Hour).
					Return(false, nil).Times(maxCodeGenerationAttempts)
			},
			expectedError: errCodeGenerationAttemptsExceeded, expectedCalls: maxCodeGenerationAttempts,
		},
		{
			name:      "get success",
			operation: "get",
			setupMockStorage: func(mockStorage *repository_mocks.URLStorage) {
				mockStorage.On("GetURL", mock.Anything, "abc1234").Return("https://example.com", nil)
			},
			expectedURL: "https://example.com",
		},
		{
			name:      "get URL does not exist",
			operation: "get",
			setupMockStorage: func(mockStorage *repository_mocks.URLStorage) {
				mockStorage.On("GetURL", mock.Anything, "abc1234").Return("", repository.ErrUrlDoesNotExist)
			},
			expectedError: ErrUrlDoesNotExist,
		},
		{
			name:      "get repository error",
			operation: "get",
			setupMockStorage: func(mockStorage *repository_mocks.URLStorage) {
				mockStorage.On("GetURL", mock.Anything, "abc1234").Return("", errors.New("storage unavailable"))
			},
			expectedError: errors.New("storage unavailable"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockStorage := &repository_mocks.URLStorage{}
			tc.setupMockStorage(mockStorage)
			shortenSvc := NewShortenURL(mockStorage)

			if tc.operation == "create" {
				code, err := shortenSvc.CreateShorternLink(context.Background(), "https://example.com", time.Hour)
				assert.ErrorIs(t, err, tc.expectedError)
				if tc.verifyCode {
					assert.Len(t, code, defaultCodeLength)
					for _, character := range code {
						assert.Contains(t, charset, string(character))
					}
				} else {
					assert.Empty(t, code)
				}
				assert.Equal(t, tc.expectedCalls, len(mockStorage.Calls))
			} else {
				url, err := shortenSvc.GetUrlFromCode(context.Background(), "abc1234")
				assert.Equal(t, tc.expectedURL, url)
				if tc.name == "get repository error" {
					assert.EqualError(t, err, tc.expectedError.Error())
				} else {
					assert.ErrorIs(t, err, tc.expectedError)
				}
			}
			mockStorage.AssertExpectations(t)
		})
	}
}
