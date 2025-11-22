package handlers

import (
	"encoding/json"
	"net/http"
	"pr-reviewer-service-AVITO2025/internal/api"
	"pr-reviewer-service-AVITO2025/internal/service"
	"pr-reviewer-service-AVITO2025/internal/storage/postgres"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) SetUserActive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		UserID   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.NewErrorResponse("VALIDATION_ERROR", "invalid json"))
		return
	}

	user, err := h.userService.SetUserActive(ctx, req.UserID, req.IsActive)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		switch err {
		case postgres.ErrUserNotFound:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(api.NewErrorResponse("NOT_FOUND", "user not found"))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.NewErrorResponse("INTERNAL_ERROR", "internal server error"))
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": user,
	})
}
