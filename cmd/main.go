package main

import (
	"context"
	"log"
	"net/http"
	"pr-reviewer-service-AVITO2025/internal/database"
	myhttp "pr-reviewer-service-AVITO2025/internal/http"
	"pr-reviewer-service-AVITO2025/internal/http/handlers"
	"pr-reviewer-service-AVITO2025/internal/service"
	"pr-reviewer-service-AVITO2025/internal/storage/postgres"
)

func main() {
	ctx := context.Background()

	databaseURL := "postgresql://postgres:pr_service_password@localhost:5432/pr_service?sslmode=disable"

	pool, err := database.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer pool.Close()

	teamRepo := postgres.NewTeamRepo(pool)
	teamService := service.NewTeamService(teamRepo)
	teamHandler := handlers.NewTeamHandler(teamService)

	router := myhttp.NewRouter(teamHandler)

	log.Println("Server start on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
