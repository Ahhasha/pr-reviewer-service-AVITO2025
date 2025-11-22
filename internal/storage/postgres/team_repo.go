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
	teamQuery := `INSERT INTO teams (id, name) 
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

		memberQuery := `INSERT INTO team_members (team_id, user_id)
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
