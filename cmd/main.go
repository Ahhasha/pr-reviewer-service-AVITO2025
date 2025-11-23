package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"pr-reviewer-service-AVITO2025/internal/config"
	"pr-reviewer-service-AVITO2025/internal/database"
	myhttp "pr-reviewer-service-AVITO2025/internal/http"
	"pr-reviewer-service-AVITO2025/internal/http/handlers"
	"pr-reviewer-service-AVITO2025/internal/service"
	"pr-reviewer-service-AVITO2025/internal/storage/postgres"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()

	lgr := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	cfg := config.Load()
	databaseURL := cfg.DatabaseURL()

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

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		lgr.Info("Server start on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()
	log.Println("timeout of 5 seconds.")

	log.Println("Server exiting")
}
