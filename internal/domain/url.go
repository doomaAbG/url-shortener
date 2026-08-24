package domain

import (
	"context"
	"time"
)

type URL struct {
	ID        int64     `json:"id"`
	Original  string    `json:"original_url"`
	Alias     string    `json:"alias"`
	CreatedAt time.Time `json:"created_at"`
}

type URLRepository interface {
	Save(ctx context.Context, url *URL) error
	GetByAlias(ctx context.Context, alias string) (*URL, error)
}
