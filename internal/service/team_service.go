package service

import (
	"context"

	"pr-reviewer-service-AVITO2025/internal/api"
	"pr-reviewer-service-AVITO2025/internal/storage/postgres"
)

type TeamService struct {
	teamRepo *postgres.TeamRepo
}

func NewTeamService(teamRepo *postgres.TeamRepo) *TeamService {
	return &TeamService{teamRepo: teamRepo}
}

func (s *TeamService) CreateTeam(ctx context.Context, team *api.Team) (*api.Team, error) {

	createdTeam, err := s.teamRepo.Create(ctx, team)
	if err != nil {
		return nil, err
	}
	return createdTeam, nil
}
