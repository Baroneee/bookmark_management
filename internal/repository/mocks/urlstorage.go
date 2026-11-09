package urlstorage_repo_mocks

import (
	"context"

	repository "github.com/Baroneee/EBVN_Assignment_1/internal/repository"
	"github.com/stretchr/testify/mock"
)

// URLStorage is a mock implementation of repository.URLStorage.
type URLStorage struct {
	mock.Mock
}

// StoreURL calls the configured mock expectation.
func (m *URLStorage) StoreURL(ctx context.Context, key, url string) error {
	return m.Called(ctx, key, url).Error(0)
}

// StoreURLIfAbsent calls the configured mock expectation.
func (m *URLStorage) StoreURLIfAbsent(ctx context.Context, key, url string) (bool, error) {
	args := m.Called(ctx, key, url)
	return args.Bool(0), args.Error(1)
}

// GetURL calls the configured mock expectation.
func (m *URLStorage) GetURL(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

var _ repository.URLStorage = (*URLStorage)(nil)
