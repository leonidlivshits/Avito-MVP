package postgres

import (
	"context"
	// "errors"
	"strings"
	"time"

	"github.com/leonidlivshits/Avito-MVP/internal/domain"
	"github.com/leonidlivshits/Avito-MVP/internal/domain/model"
	"github.com/leonidlivshits/Avito-MVP/internal/ports"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPRRepo struct {
	pool *pgxpool.Pool
}

var _ ports.PRRepository = (*PostgresPRRepo)(nil)

func NewPostgresPRRepo(pool *pgxpool.Pool) *PostgresPRRepo {
	return &PostgresPRRepo{pool: pool}
}

func (r *PostgresPRRepo) GetByID(prID string) (*model.PullRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.pool.QueryRow(ctx, `
		SELECT pull_request_id, pull_request_name, author_id, status, need_more_reviewers, created_at, merged_at
		FROM pull_requests WHERE pull_request_id = $1
	`, prID)

	var id, name, authorID, statusText string
	var needMore bool
	var createdAt time.Time
	var mergedAt *time.Time
	if err := row.Scan(&id, &name, &authorID, &statusText, &needMore, &createdAt, &mergedAt); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	rows, err := r.pool.Query(ctx, `SELECT user_id FROM pr_reviewers WHERE pr_id = $1 ORDER BY assigned_at`, prID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var revs []string
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		revs = append(revs, uid)
	}

	pr := &model.PullRequest{
		PullRequestID:     id,
		PullRequestName:   name,
		AuthorID:          authorID,
		Status:            model.PRStatus(statusText),
		AssignedReviewers: revs,
		NeedMoreReviewers: needMore,
		CreatedAt:         createdAt,
		MergedAt:          mergedAt,
	}
	return pr, nil
}

func (r *PostgresPRRepo) Create(pr *model.PullRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var exists bool
	row := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1)`, pr.PullRequestID)
	if err := row.Scan(&exists); err != nil {
		return err
	}
	if exists {
		return domain.ErrPRExists
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, need_more_reviewers, created_at, merged_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, pr.PullRequestID, pr.PullRequestName, pr.AuthorID, string(pr.Status), pr.NeedMoreReviewers, pr.CreatedAt, pr.MergedAt)
	if err != nil {
		return err
	}

	for _, uid := range pr.AssignedReviewers {
		if _, err := tx.Exec(ctx, `
			INSERT INTO pr_reviewers (pr_id, user_id, assigned_at) VALUES ($1, $2, $3)
		`, pr.PullRequestID, uid, time.Now().UTC()); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (r *PostgresPRRepo) Update(pr *model.PullRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	row := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1)`, pr.PullRequestID)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return domain.ErrNotFound
	}

	_, err = tx.Exec(ctx, `
		UPDATE pull_requests
		SET pull_request_name = $1, author_id = $2, status = $3, need_more_reviewers = $4, created_at = $5, merged_at = $6
		WHERE pull_request_id = $7
	`, pr.PullRequestName, pr.AuthorID, string(pr.Status), pr.NeedMoreReviewers, pr.CreatedAt, pr.MergedAt, pr.PullRequestID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM pr_reviewers WHERE pr_id = $1`, pr.PullRequestID)
	if err != nil {
		return err
	}
	for _, uid := range pr.AssignedReviewers {
		if _, err := tx.Exec(ctx, `INSERT INTO pr_reviewers (pr_id, user_id, assigned_at) VALUES ($1, $2, $3)`, pr.PullRequestID, uid, time.Now().UTC()); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (r *PostgresPRRepo) ListByReviewer(userID string) ([]*model.PullRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.pool.Query(ctx, `
		SELECT p.pull_request_id, p.pull_request_name, p.author_id, p.status, p.need_more_reviewers, p.created_at, p.merged_at
		FROM pull_requests p
		JOIN pr_reviewers r ON r.pr_id = p.pull_request_id
		WHERE r.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*model.PullRequest
	for rows.Next() {
		var id, name, authorID, statusText string
		var needMore bool
		var createdAt time.Time
		var mergedAt *time.Time
		if err := rows.Scan(&id, &name, &authorID, &statusText, &needMore, &createdAt, &mergedAt); err != nil {
			return nil, err
		}
		rvs := make([]string, 0)
		rrows, err := r.pool.Query(ctx, `SELECT user_id FROM pr_reviewers WHERE pr_id = $1 ORDER BY assigned_at`, id)
		if err != nil {
			return nil, err
		}
		for rrows.Next() {
			var uid string
			if err := rrows.Scan(&uid); err != nil {
				rrows.Close()
				return nil, err
			}
			rvs = append(rvs, uid)
		}
		rrows.Close()

		out = append(out, &model.PullRequest{
			PullRequestID:     id,
			PullRequestName:   name,
			AuthorID:          authorID,
			Status:            model.PRStatus(statusText),
			AssignedReviewers: rvs,
			NeedMoreReviewers: needMore,
			CreatedAt:         createdAt,
			MergedAt:          mergedAt,
		})
	}
	return out, nil
}

func (r *PostgresPRRepo) ListOpen() ([]*model.PullRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.pool.Query(ctx, `
		SELECT pull_request_id, pull_request_name, author_id, status, need_more_reviewers, created_at, merged_at
		FROM pull_requests WHERE status = 'OPEN'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*model.PullRequest
	for rows.Next() {
		var id, name, authorID, statusText string
		var needMore bool
		var createdAt time.Time
		var mergedAt *time.Time
		if err := rows.Scan(&id, &name, &authorID, &statusText, &needMore, &createdAt, &mergedAt); err != nil {
			return nil, err
		}
		rvs := make([]string, 0)
		rrows, err := r.pool.Query(ctx, `SELECT user_id FROM pr_reviewers WHERE pr_id = $1 ORDER BY assigned_at`, id)
		if err != nil {
			return nil, err
		}
		for rrows.Next() {
			var uid string
			if err := rrows.Scan(&uid); err != nil {
				rrows.Close()
				return nil, err
			}
			rvs = append(rvs, uid)
		}
		rrows.Close()

		out = append(out, &model.PullRequest{
			PullRequestID:     id,
			PullRequestName:   name,
			AuthorID:          authorID,
			Status:            model.PRStatus(statusText),
			AssignedReviewers: rvs,
			NeedMoreReviewers: needMore,
			CreatedAt:         createdAt,
			MergedAt:          mergedAt,
		})
	}
	return out, nil
}
