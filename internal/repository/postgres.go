package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/doomaAbG/url-shortener/internal/domain"
)

type PostgresRepo struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, connString string) (*PostgresRepo, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &PostgresRepo{pool: pool}, nil
}

func (r *PostgresRepo) Save(ctx context.Context, url *domain.URL) error {
	query := `INSERT INTO urls (original_url, alias, created_at) VALUES ($1, $2, $3) RETURNING id`
	err := r.pool.QueryRow(ctx, query, url.Original, url.Alias, url.CreatedAt).Scan(&url.ID)
	if err != nil {
		return fmt.Errorf("failed to save url: %w", err)
	}
	return nil
}

func (r *PostgresRepo) GetByAlias(ctx context.Context, alias string) (*domain.URL, error) {
	query := `SELECT id, original_url, alias, created_at FROM urls WHERE alias = $1`

	var url domain.URL
	err := r.pool.QueryRow(ctx, query, alias).Scan(&url.ID, &url.Original, &url.Alias, &url.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get url: %w", err)
	}

	return &url, nil
}

func (r *PostgresRepo) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}
