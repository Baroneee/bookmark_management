package service

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/Baroneee/EBVN_Assignment_1/internal/repository"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const defaultCodeLength = 7
const maxCodeGenerationAttempts = 10

var errCodeGenerationAttemptsExceeded = errors.New("unable to generate a unique short URL code")

//go:generate mockery --name=ShortenURL --filename=shorternurl.go --outpkg svc_mocks
// ShortenURL creates and stores shortened URL codes.
type ShortenURL interface {
	CreateShorternLink(ctx context.Context, url string) (string, error)
}

type shortenURL struct {
	r repository.URLStorage
}

// NewShortenURL creates a URL-shortening service backed by r.
func NewShortenURL(r repository.URLStorage) ShortenURL {
	return &shortenURL{r: r}
}

// CreateShorternLink generates a random code, stores the URL under that code,
// and returns the generated code.
func (s *shortenURL) CreateShorternLink(ctx context.Context, url string) (string, error) {
	for attempt := 0; attempt < maxCodeGenerationAttempts; attempt++ {
		code, err := generateCode()
		if err != nil {
			return "", err
		}

		created, err := s.r.StoreURLIfAbsent(ctx, code, url)
		if err != nil {
			return "", err
		}
		if created {
			return code, nil
		}
	}

	return "", errCodeGenerationAttemptsExceeded
}

func generateCode() (string, error) {
	code := make([]byte, defaultCodeLength)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[n.Int64()]
	}
	return string(code), nil
}
