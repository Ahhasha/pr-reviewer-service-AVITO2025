package service

import (
	"context"
	"fmt"
	"log/slog"
	"pr-reviewer-service-AVITO2025/internal/storage/postgres"
)

type PRService struct {
	prRepo   *postgres.PullRequestRepo
	teamRepo *postgres.TeamRepo
	lgr      *slog.Logger
}

func NewPRService(prRepo *postgres.PullRequestRepo, teamRepo *postgres.TeamRepo, lgr *slog.Logger) *PRService {
	return &PRService{
		prRepo:   prRepo,
		teamRepo: teamRepo,
		lgr:      lgr,
	}
}

func (s *PRService) FindReviewersForAuthor(ctx context.Context, authorID string) ([]string, error) {
	s.lgr.Info("finding reviewers for author", "author_id", authorID)

	team, err := s.teamRepo.GetByUserID(ctx, authorID)
	if err != nil {
		s.lgr.Error("failed to find author team", "author_id", authorID, "error", err)
		return nil, fmt.Errorf("author has no team: %w", err)
	}

	var candidates []string
	for _, member := range team.Members {
		if member.UserId != authorID && member.IsActive {
			candidates = append(candidates, member.UserId)
		}
	}

	s.lgr.Info("found candidates", "author_id", authorID, "candidates_count", len(candidates), "candidates", candidates)

	var reviewers []string
	if len(candidates) > 2 {
		reviewers = candidates[:2]
	} else {
		reviewers = candidates
	}

	s.lgr.Info("assigned reviewers", "author_id", authorID, "reviewers", reviewers)

	return reviewers, nil
}
