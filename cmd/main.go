package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"pr-reviewer-service-AVITO2025/internal/database"
	myhttp "pr-reviewer-service-AVITO2025/internal/http"
	"pr-reviewer-service-AVITO2025/internal/http/handlers"
	"pr-reviewer-service-AVITO2025/internal/service"
	"pr-reviewer-service-AVITO2025/internal/storage/postgres"
)

func main() {
	ctx := context.Background()

	lgr := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	databaseURL := "postgresql://postgres:pr_service_password@localhost:5432/pr_service?sslmode=disable"

	pool, err := database.NewPool(ctx, databaseURL)
	if err != nil {
		lgr.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	teamRepo := postgres.NewTeamRepo(pool)
	userRepo := postgres.NewUserRepo(pool)
	prRepo := postgres.NewPullRequestRepo(pool)

	teamService := service.NewTeamService(teamRepo, lgr)
	userService := service.NewUserService(userRepo, lgr)
	prService := service.NewPRService(prRepo, teamRepo, lgr)

	teamHandler := handlers.NewTeamHandler(teamService)
	userHandler := handlers.NewUserHandler(userService)
	prHandler := handlers.NewPRHandler(prService, prRepo)

	router := myhttp.NewRouter(teamHandler, userHandler, prHandler)

	lgr.Info("Server start on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
