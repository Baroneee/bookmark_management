package service

import (
	"testing"

	repository_mocks "github.com/Baroneee/EBVN_Assignment_1/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestShortenURL_CreateShorternLink verifies code generation, persistence, and
// retry behavior when a generated code already exists.
func TestShortenURL_CreateShorternLink(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		setupMockStorage func(*repository_mocks.URLStorage)
		expectedError    error
		verifyCode       bool
		expectedCalls    int
	}{
		{
			name: "success",
			setupMockStorage: func(mockStorage *repository_mocks.URLStorage) {
				mockStorage.On("StoreURLIfAbsent", mock.Anything, mock.Anything, "https://example.com").
					Return(true, nil).Once()
			},
			verifyCode:    true,
			expectedCalls: 1,
		},
		{
			name: "retry after collision",
			setupMockStorage: func(mockStorage *repository_mocks.URLStorage) {
				mockStorage.On("StoreURLIfAbsent", mock.Anything, mock.Anything, "https://example.com").
					Return(false, nil).Once()
				mockStorage.On("StoreURLIfAbsent", mock.Anything, mock.Anything, "https://example.com").
					Return(true, nil).Once()
			},
			verifyCode:    true,
			expectedCalls: 2,
		},
		{
			name: "storage error",
			setupMockStorage: func(mockStorage *repository_mocks.URLStorage) {
				mockStorage.On("StoreURLIfAbsent", mock.Anything, mock.Anything, "https://example.com").
					Return(false, assert.AnError).Once()
			},
			expectedError: assert.AnError,
			expectedCalls: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockStorage := &repository_mocks.URLStorage{}
			tc.setupMockStorage(mockStorage)
			shortenSvc := NewShortenURL(mockStorage)

			code, err := shortenSvc.CreateShorternLink(ctx, "https://example.com")

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
			mockStorage.AssertExpectations(t)
		})
	}
}
