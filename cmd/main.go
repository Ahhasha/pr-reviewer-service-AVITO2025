package main

import (
	"context"
	"log"
	"net/http"
	"pr-reviewer-service-AVITO2025/internal/database"

	"github.com/go-chi/chi/v5"
)

func main() {
	ctx := context.Background()

	databaseURL := "postgresql://postgres:pr_service_password@localhost:5432/pr_service?sslmode=disable"

	pool, err := database.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer pool.Close()

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	log.Println("Server start on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
