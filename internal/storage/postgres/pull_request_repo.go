package postgres

import (
	"context"
	"errors"
	"fmt"
	"pr-reviewer-service-AVITO2025/internal/api"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrPRNotFound = errors.New("pull request not found")
var ErrPRExists = errors.New("pull request already exists")
var ErrPRMerged = errors.New("pull request is merged")
var ErrNotAssigned = errors.New("reviewer is not assigned to this PR")
var ErrNoCandidate = errors.New("no active replacement candidate in team")

type PullRequestRepo struct {
	pool *pgxpool.Pool
}

func NewPullRequestRepo(pool *pgxpool.Pool) *PullRequestRepo {
	return &PullRequestRepo{
		pool: pool,
	}
}

func (r *PullRequestRepo) Create(ctx context.Context, pr *api.PullRequest, reviewers []string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	prQuery := `
			INSERT INTO pull_requests (id, name, author_id, status)
			VALUES ($1, $2, $3, $4)
	`
	_, err = tx.Exec(ctx, prQuery, pr.PullRequestId, pr.PullRequestName, pr.AuthorId, "OPEN")
	if err != nil {
		return ErrPRExists
	}

	for _, reviewerID := range reviewers {
		reviewerQuery := `
					INSERT INTO pr_reviewers (pr_id, user_id)
					VALUES ($1, $2)
		`
		_, err := tx.Exec(ctx, reviewerQuery, pr.PullRequestId, reviewerID)
		if err != nil {
			return fmt.Errorf("add reviewer: %w", err)
		}
	}
	return tx.Commit(ctx)
}

func (r *PullRequestRepo) GetByID(ctx context.Context, prID string) (*api.PullRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var pr api.PullRequest
	var createdAt, mergedAt *time.Time

	query := `
			SELECT id, name, author_id, status, created_at, merged_at
			FROM pull_request
			WHERE id = $1
	`

	err := r.pool.QueryRow(ctx, query, prID).Scan(&pr.PullRequestId, &pr.PullRequestName, &pr.AuthorId, &pr.Status, &createdAt, &mergedAt)
	if err != nil {
		return nil, ErrPRNotFound
	}

	if createdAt != nil {
		pr.CreatedAt = createdAt
	}

	if mergedAt != nil {
		pr.MergedAt = mergedAt
	}

	reviewersQuery := `
					SELECT user_id
					FROM pr_reviewers
					WHERE pr_id = $1
	`
	rows, err := r.pool.Query(ctx, reviewersQuery, prID)
	if err != nil {
		return nil, fmt.Errorf("get reviewers: %w", err)
	}
	defer rows.Close()

	pr.AssignedReviewers = []string{}

	for rows.Next() {
		var reviewerID string
		err := rows.Scan(&reviewerID)
		if err != nil {
			return nil, fmt.Errorf("scan reviwer: %w", err)
		}

		pr.AssignedReviewers = append(pr.AssignedReviewers, reviewerID)
	}
	return &pr, nil
}

func (r *PullRequestRepo) Merge(ctx context.Context, prID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
			UPDATE pull_requests
			SET status = 'MERGED', merged_at = NOW()
			WHERE id = $1 AND status = 'OPEN'
	`

	result, err := r.pool.Exec(ctx, query, prID)
	if err != nil {
		return fmt.Errorf("merge PR: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrPRNotFound
	}
	return nil
}
