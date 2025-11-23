package validation

import (
	"pr-reviewer-service-AVITO2025/internal/api"
)

type ValidationError struct {
	Place   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func ValidateCreateTeam(team *api.Team) error {
	if team.TeamName == "" {
		return &ValidationError{
			Place:   "team_name",
			Message: "Cannot be empty",
		}
	}

	if len(team.TeamName) < 2 || len(team.TeamName) > 50 {
		return &ValidationError{
			Place:   "team_name",
			Message: "Team name must be 1..50 chars",
		}
	}

	if len(team.Members) == 0 {
		return &ValidationError{
			Place:   "members",
			Message: "The team must contain at least one member",
		}
	}

	for _, member := range team.Members {
		if member.UserId == "" {
			return &ValidationError{
				Place:   "members.user_id",
				Message: "The user must have a user ID.",
			}
		}

		if member.Username == "" {
			return &ValidationError{
				Place:   "members.username",
				Message: "Cannot be empty",
			}
		}

		if len(member.Username) < 1 || len(member.Username) > 100 {
			return &ValidationError{
				Place:   "members.username",
				Message: "User name must be 1..100 chars",
			}
		}
	}
	return nil
}

func ValidateCreatePR(authorID, prID, name string) error {
	if authorID == "" {
		return &ValidationError{
			Place:   "author_id",
			Message: "Cannot be empty",
		}
	}

	if prID == "" {
		return &ValidationError{
			Place:   "pull_request_id",
			Message: "Cannot be empty",
		}
	}

	if name == "" {
		return &ValidationError{
			Place:   "pull_request_name",
			Message: "Cannot be empty",
		}
	}

	if len(name) < 1 || len(name) > 200 {
		return &ValidationError{
			Place:   "pull_request_name",
			Message: "PR name must be 1..200 chars",
		}
	}

	return nil
}
