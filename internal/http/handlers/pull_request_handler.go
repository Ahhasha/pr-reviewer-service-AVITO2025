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
		json.NewEncoder(w).Encode(api.NewErrorResponse("METHOD_NOT_ALLOWED", "Only POST method"))
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

func (h *PRHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := r.URL.Query().Get("user_id")
	if err := validation.ValidateUserID(userID); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.NewErrorResponse("VALIDATION_ERROR", "user_id cannot be empty"))
		return
	}

	prs, err := h.prRepo.GetByReviewer(ctx, userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.NewErrorResponse("INTERNAL_ERROR", "Failed to get user PR"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":       userID,
		"pull_requests": prs,
	})
}

func (h *PRHandler) MergePR(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(api.NewErrorResponse("METHOD_NOT_ALLOWED", "Only POST method"))
		return
	}

	var req struct {
		PRID string `json:"pull_request_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.NewErrorResponse("VALIDATION_ERROR", "Invalid JSON format"))
		return
	}

	if err := validation.ValidateMergePR(req.PRID); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	err := h.prRepo.Merge(ctx, req.PRID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		switch err {
		case postgres.ErrPRNotFound:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(api.NewErrorResponse("NOT_FOUND", "Pull request not found"))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.NewErrorResponse("INTERNAL_ERROR", "Failed to merge PR"))
		}
		return
	}

	pr, err := h.prRepo.GetByID(ctx, req.PRID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.NewErrorResponse("INTERNAL_ERROR", "Failed to get updated PR"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pr": pr,
	})
}

func (h *PRHandler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(api.NewErrorResponse("METHOD_NOT_ALLOWED", "Only POST method"))
		return
	}

	var req struct {
		PRID      string `json:"pull_request_id"`
		OldUserID string `json:"old_user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.NewErrorResponse("VALIDATION_ERROR", "Invalid JSON format"))
		return
	}

	if err := validation.ValidateReassignPR(req.PRID, req.OldUserID); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	newUserID, err := h.prService.ReassignReviewer(ctx, req.PRID, req.OldUserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		switch err {
		case postgres.ErrPRNotFound:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(api.NewErrorResponse("NOT_FOUND", "Pull request not found"))
		case postgres.ErrNoCandidate:
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(api.NewErrorResponse("NO_CANDIDATE", "No active replacement candidate in team"))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.NewErrorResponse("INTERNAL_ERROR", "Failed to find replacement"))
		}
		return
	}

	err = h.prRepo.Reassign(ctx, req.PRID, req.OldUserID, newUserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		switch err {
		case postgres.ErrPRNotFound:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(api.NewErrorResponse("NOT_FOUND", "Pull request not found"))
		case postgres.ErrPRMerged:
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(api.NewErrorResponse("PR_MERGED", "Cannot reassign on merged PR"))
		case postgres.ErrNotAssigned:
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(api.NewErrorResponse("NOT_ASSIGNED", "Reviewer is not assigned to this PR"))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.NewErrorResponse("INTERNAL_ERROR", "Failed to reassign reviewer"))
		}
		return
	}

	pr, err := h.prRepo.GetByID(ctx, req.PRID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.NewErrorResponse("INTERNAL_ERROR", "Failed to get updated PR"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pr":          pr,
		"replaced_by": newUserID,
	})
}
