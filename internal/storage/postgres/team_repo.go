package postgres

import (
	"context"
	"errors"
	"fmt"
	"pr-reviewer-service-AVITO2025/internal/api"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrTeamNotFound = errors.New("team not found")
var ErrTeamExists = errors.New("team already exists")

type TeamRepo struct {
	pool *pgxpool.Pool
}

func NewTeamRepo(pool *pgxpool.Pool) *TeamRepo {
	return &TeamRepo{pool: pool}
}

func (r *TeamRepo) Create(ctx context.Context, team *api.Team) (*api.Team, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	teamID := uuid.New().String()
	teamQuery := `
				INSERT INTO teams (id, team_name) 
				VALUES ($1, $2)
	`

	_, err = tx.Exec(ctx, teamQuery, teamID, team.TeamName)
	if err != nil {
		return nil, ErrTeamExists
	}

	for _, member := range team.Members {
		userQuery := `
					INSERT INTO users (id, username, is_active)
					VALUES ($1, $2, $3)
					ON CONFLICT (id) DO UPDATE SET
					username = EXCLUDED.username,
					is_active = EXCLUDED.is_active
		`
		_, err = tx.Exec(ctx, userQuery, member.UserId, member.Username, member.IsActive)
		if err != nil {
			return nil, fmt.Errorf("create or update user: %w", err)
		}

		memberQuery := `
						INSERT INTO team_members (team_id, user_id)
						VALUES ($1, $2)
		`
		_, err = tx.Exec(ctx, memberQuery, teamID, member.UserId)
		if err != nil {
			return nil, fmt.Errorf("add user in team: %w", err)
		}
	}

	err = tx.Commit(ctx)

	if err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return team, nil
}

func (r *TeamRepo) GetByName(ctx context.Context, teamName string) (*api.Team, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var teamID string
	err := r.pool.QueryRow(ctx, "SELECT id FROM teams WHERE team_name = $1", teamName).Scan(&teamID)
	if err != nil {
		return nil, ErrTeamNotFound
	}

	query := `
			SELECT u.id, u.username, u.is_active
			FROM users u
			JOIN team_members tm ON u.id = tm.user_id
			WHERE tm.team_id = $1
	`
	rows, err := r.pool.Query(ctx, query, teamID)
	if err != nil {
		return nil, fmt.Errorf("query team members: %w", err)
	}
	defer rows.Close()

	var team api.Team
	team.TeamName = teamName
	team.Members = []api.TeamMember{}

	for rows.Next() {
		var member api.TeamMember
		err := rows.Scan(&member.UserId, &member.Username, &member.IsActive)
		if err != nil {
			return nil, fmt.Errorf("scan team member: %w", err)
		}
		team.Members = append(team.Members, member)
	}
	return &team, nil
}
