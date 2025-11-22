package service

import (
	"context"
	"pr-reviewer-service-AVITO2025/internal/api"
	"pr-reviewer-service-AVITO2025/internal/storage/postgres"
)

type UserService struct {
	userRepo *postgres.UserRepo
}

func NewUserService(userRepo *postgres.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) SetUserActive(ctx context.Context, userID string, isActive bool) (*api.User, error) {
	user, err := s.userRepo.SetUserActive(ctx, userID, isActive)
	if err != nil {
		return nil, err
	}
	return user, nil
}
