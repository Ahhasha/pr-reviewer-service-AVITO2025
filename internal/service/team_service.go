package service

import (
	"context"
	"log/slog"

	"pr-reviewer-service-AVITO2025/internal/api"
	"pr-reviewer-service-AVITO2025/internal/storage/postgres"
)

type TeamService struct {
	teamRepo *postgres.TeamRepo
	lgr      *slog.Logger
}

func NewTeamService(teamRepo *postgres.TeamRepo, lgr *slog.Logger) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
		lgr:      lgr,
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, team *api.Team) (*api.Team, error) {

	createdTeam, err := s.teamRepo.Create(ctx, team)
	if err != nil {
		s.lgr.Error("failed to create team", "method", "CreateTeam", "team_name", team.TeamName, "error", err.Error())
		return nil, err
	}
	s.lgr.Info("team created", "team_name", team.TeamName, "members_count", len(team.Members))
	return createdTeam, nil
}

func (s *TeamService) GetTeam(ctx context.Context, teamName string) (*api.Team, error) {
	team, err := s.teamRepo.GetByName(ctx, teamName)

	if err != nil {
		s.lgr.Error("failed to get team", "method", "GetTeam", "team_name", teamName, "error", err.Error())
		return nil, err
	}
	s.lgr.Info("the team has been received", "team_name", teamName, "members_count", len(team.Members))
	return team, nil
}
