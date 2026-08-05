package main

import (
	"log"
	"net/http"

	"github.com/doomAbG/url-shortener/internal/handler"
)

func main() {
	http.HandleFunc("/health", handler.HealthCheck)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
