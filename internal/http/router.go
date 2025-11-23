package http

import (
	"net/http"
	"pr-reviewer-service-AVITO2025/internal/http/handlers"

	"github.com/go-chi/chi/v5"
)

func NewRouter(teamHandler *handlers.TeamHandler, userHandler *handlers.UserHandler, prHandler *handlers.PRHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/team/add", teamHandler.CreateTeam)
	r.Get("/team/get", teamHandler.GetTeam)

	r.Post("/users/setIsActive", userHandler.SetUserActive)
	r.Get("/users/getReview", prHandler.GetReview)

	r.Post("/pullRequest/create", prHandler.CreatePR)
	r.Post("/pullRequest/merge", prHandler.MergePR)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	return r
}
