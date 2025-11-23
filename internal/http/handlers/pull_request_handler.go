package handlers

import (
	"encoding/json"
	"net/http"
	"pr-reviewer-service-AVITO2025/internal/api"
	"pr-reviewer-service-AVITO2025/internal/service"
	"pr-reviewer-service-AVITO2025/internal/storage/postgres"
	"pr-reviewer-service-AVITO2025/internal/validation"
)

type PRHandler struct {
	prService *service.PRService
	prRepo    *postgres.PullRequestRepo
}

func NewPRHandler(prService *service.PRService, prRepo *postgres.PullRequestRepo) *PRHandler {
	return &PRHandler{
		prService: prService,
		prRepo:    prRepo,
	}
}

func (h *PRHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(api.NewErrorResponse("METHOD_NOT_ALLOWED", "Only POST method is allowed"))
		return
	}

	var req struct {
		AuthorID string `json:"author_id"`
		PRID     string `json:"pull_request_id"`
		Name     string `json:"pull_request_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.NewErrorResponse("VALIDATION_ERROR", "Invalid JSON format"))
		return
	}

	if err := validation.ValidateCreatePR(req.AuthorID, req.PRID, req.Name); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	reviewers, err := h.prService.FindReviewersForAuthor(ctx, req.AuthorID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.NewErrorResponse("NOT_FOUND", err.Error()))
		return
	}

	pr := &api.PullRequest{
		PullRequestId:     req.PRID,
		PullRequestName:   req.Name,
		AuthorId:          req.AuthorID,
		Status:            api.PullRequestStatusOPEN,
		AssignedReviewers: reviewers,
	}

	err = h.prRepo.Create(ctx, pr, reviewers)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.NewErrorResponse("INTERNAL_ERROR", "Failed to create PR"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pr": pr,
	})
}
