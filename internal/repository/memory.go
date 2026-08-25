package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/doomaAbG/url-shortener/internal/domain"
)

var ErrNotFound = errors.New("url not found")

type MemoryRepo struct {
	mu   sync.RWMutex
	urls map[string]*domain.URL
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		urls: make(map[string]*domain.URL),
	}
}

func (r *MemoryRepo) Save(ctx context.Context, url *domain.URL) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.urls[url.Alias] = url
	return nil
}

func (r *MemoryRepo) GetByAlias(ctx context.Context, alias string) (*domain.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, exists := r.urls[alias]
	if !exists {
		return nil, ErrNotFound
	}
	return url, nil
}
