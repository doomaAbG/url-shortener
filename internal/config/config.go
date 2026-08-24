package config

import (
	"os"
)

type Config struct {
	Env   string
	Port  string
	DBURL string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:secret@127.0.0.1:5439/url_shortener?sslmode=disable"
	}

	return &Config{
		Env:   env,
		Port:  port,
		DBURL: dbURL,
	}
}
