package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/doomAbG/url-shortener/internal/config"
	"github.com/doomAbG/url-shortener/internal/handler"
	"github.com/doomAbG/url-shortener/internal/repository"
	"github.com/doomAbG/url-shortener/internal/service"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	// Инициализируем Postgres репозиторий вместо MemoryRepo
	repo, err := repository.NewPostgresRepo(ctx, cfg.DBURL)
	if err != nil {
		log.Fatalf("Failed to initialize postgres repo: %v", err)
	}
	defer repo.Close()

	urlService := service.NewURLService(repo)
	urlHandler := handler.NewURLHandler(urlService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", handler.HealthCheck)
	r.Post("/shorten", urlHandler.Shorten)
	r.Get("/{alias}", urlHandler.Redirect)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s (env: %s)...", addr, cfg.Env)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
