package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/doomAbG/url-shortener/internal/domain"
)

type PostgresRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresRepo(ctx context.Context, connString string) (*PostgresRepo, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	// Создаем таблицу, если ее нет
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		original_url TEXT NOT NULL,
		alias VARCHAR(20) UNIQUE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_urls_alias ON urls(alias);
	`

	_, err = pool.Exec(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
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
	r.pool.Close()
}
