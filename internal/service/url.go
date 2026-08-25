package service

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"github.com/doomaAbG/url-shortener/internal/domain"
)

type URLService struct {
	repo domain.URLRepository
}

func NewURLService(repo domain.URLRepository) *URLService {
	return &URLService{repo: repo}
}

func (s *URLService) Shorten(ctx context.Context, originalURL string) (*domain.URL, error) {
	alias, err := generateRandomAlias(6)
	if err != nil {
		return nil, err
	}

	url := &domain.URL{
		Original:  originalURL,
		Alias:     alias,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Save(ctx, url); err != nil {
		return nil, err
	}

	return url, nil
}

func (s *URLService) GetOriginal(ctx context.Context, alias string) (string, error) {
	url, err := s.repo.GetByAlias(ctx, alias)
	if err != nil {
		return "", err
	}
	return url.Original, nil
}

func generateRandomAlias(length int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		result[i] = chars[num.Int64()]
	}
	return string(result), nil
}
