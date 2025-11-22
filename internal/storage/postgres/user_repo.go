package postgres

import (
	"context"
	"errors"
	"pr-reviewer-service-AVITO2025/internal/api"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) SetUserActive(ctx context.Context, userID string, isActive bool) (*api.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user api.User
	user.UserId = userID

	query := `
	        UPDATE users
	        SET is_active = $1
			WHERE id = $2
			RETURNING id, username, is_active
	`

	err := r.pool.QueryRow(ctx, query, isActive, userID).Scan(&user.UserId, &user.Username, &user.IsActive)
	if err != nil {
		return nil, ErrUserNotFound
	}

	teamQuery := `
				SELECT t.team_name
				FROM teams t JOIN team_members tm ON t.id = tm.team_id
				WHERE tm.user_id = $1
	`
	err = r.pool.QueryRow(ctx, teamQuery, userID).Scan(&user.TeamName)
	if err != nil {
		user.TeamName = ""
	}
	return &user, nil
}
