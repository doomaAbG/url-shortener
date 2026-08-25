package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/doomaAbG/url-shortener/internal/config"
	"github.com/doomaAbG/url-shortener/internal/handler"
	"github.com/doomaAbG/url-shortener/internal/repository"
	"github.com/doomaAbG/url-shortener/internal/service"
)

func main() {
	// 1. Инициализируем структурированный JSON-логгер
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	// 2. Загружаем конфигурацию
	cfg := config.Load()

	// 3. Запускаем миграции БД
	if err := repository.RunMigrations(cfg.DBURL); err != nil {
		slog.Error("Migration failed", "error", err)
		os.Exit(1)
	}

	// 4. Подключаемся к PostgreSQL
	repo, err := repository.NewPostgres(context.Background(), cfg.DBURL)
	if err != nil {
		slog.Error("Failed to init postgres repo", "error", err)
		os.Exit(1)
	}
	defer repo.Close()

	// 5. Собираем слои приложения
	svc := service.NewURLService(repo)
	h := handler.NewURLHandler(svc)

	// 6. Настраиваем роутер и Middleware
	r := chi.NewRouter()

	r.Use(middleware.RequestID) // Уникальный ID для каждого запроса
	r.Use(middleware.RealIP)    // Настоящий IP пользователя
	r.Use(middleware.Logger)    // Логирование HTTP-запросов
	r.Use(middleware.Recoverer) // Защита от падений (recover от panic)

	r.Post("/shorten", h.Shorten)
	r.Get("/{alias}", h.Redirect)

	slog.Info("Server started", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		slog.Error("Server stopped with error", "error", err)
	}
}
