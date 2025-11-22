package service

import (
	"context"
	"errors"
	"fmt"

	"pr-reviewer-service-AVITO2025/internal/api"
	"pr-reviewer-service-AVITO2025/internal/storage/postgres"
	"pr-reviewer-service-AVITO2025/internal/validation"
)

type TeamService struct {
	teamRepo *postgres.TeamRepo
}

func NewTeamService(teamRepo *postgres.TeamRepo) *TeamService {
	return &TeamService{teamRepo: teamRepo}
}

func (s *TeamService) CreateTeam(ctx context.Context, team *api.Team) (*api.Team, error) {
	if err := validation.ValidateCreateTeam(team); err != nil {
		return nil, err
	}

	createdTeam, err := s.teamRepo.Create(ctx, team)
	if err != nil {
		if errors.Is(err, postgres.ErrTeamExists) {
			return nil, &validation.ValidationError{
				Place:   "team_name",
				Message: "team already exists",
			}
		}
		return nil, fmt.Errorf("create team: %w", err)
	}
	return createdTeam, nil
}
