package handlers

import (
	"encoding/json"
	"net/http"

	"pr-reviewer-service-AVITO2025/internal/api"
	"pr-reviewer-service-AVITO2025/internal/service"
	"pr-reviewer-service-AVITO2025/internal/storage/postgres"
	"pr-reviewer-service-AVITO2025/internal/validation"
)

type TeamHandler struct {
	teamService *service.TeamService
}

func NewTeamHandler(teamService *service.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req api.Team
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.NewErrorResponse("VALIDATION_ERROR", "Invalid JSON format"))
		return
	}

	if err := validation.ValidateCreateTeam(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}
	team, err := h.teamService.CreateTeam(ctx, &req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		switch err {
		case postgres.ErrTeamExists:
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(api.NewErrorResponse("TEAM_EXISTS", err.Error()))
		case postgres.ErrTeamNotFound:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(api.NewErrorResponse("NOT_FOUND", err.Error()))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.NewErrorResponse("INTERNAL_ERROR", "Internal server error"))
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"team": team,
	})
}
