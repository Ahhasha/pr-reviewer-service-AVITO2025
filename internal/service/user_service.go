package service

import (
	"context"
	"log/slog"
	"pr-reviewer-service-AVITO2025/internal/api"
	"pr-reviewer-service-AVITO2025/internal/storage/postgres"
)

type UserService struct {
	userRepo *postgres.UserRepo
	lgr      *slog.Logger
}

func NewUserService(userRepo *postgres.UserRepo, lgr *slog.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		lgr:      lgr,
	}
}

func (s *UserService) SetUserActive(ctx context.Context, userID string, isActive bool) (*api.User, error) {
	user, err := s.userRepo.SetUserActive(ctx, userID, isActive)
	if err != nil {
		s.lgr.Error("failed to set user active", "method", "SetUserActive", "user_id", userID, "is_active", isActive, "error", err)
		return nil, err
	}
	s.lgr.Info("user active update", "user_id", userID, "is_active", isActive, "username", user.Username)
	return user, nil
}
