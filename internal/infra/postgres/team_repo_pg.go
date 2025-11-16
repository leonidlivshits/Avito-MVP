package postgres

import (
	"context"
	"time"

	"github.com/leonidlivshits/Avito-MVP/internal/domain"
	"github.com/leonidlivshits/Avito-MVP/internal/domain/model"
	"github.com/leonidlivshits/Avito-MVP/internal/ports"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresTeamRepo struct {
	pool     *pgxpool.Pool
	userRepo ports.UserRepository
}

var _ ports.TeamRepository = (*PostgresTeamRepo)(nil)

func NewPostgresTeamRepo(pool *pgxpool.Pool, ur ports.UserRepository) *PostgresTeamRepo {
	return &PostgresTeamRepo{pool: pool, userRepo: ur}
}

func (r *PostgresTeamRepo) GetByName(teamName string) (*model.Team, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.pool.QueryRow(ctx, `SELECT team_name FROM teams WHERE team_name = $1`, teamName)
	var tn string
	if err := row.Scan(&tn); err != nil {
		return nil, domain.ErrNotFound
	}

	rows, err := r.pool.Query(ctx, `
		SELECT user_id, username, is_active FROM users WHERE team_name = $1
	`, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]model.TeamMember, 0)
	for rows.Next() {
		var uid, username string
		var isActive bool
		if err := rows.Scan(&uid, &username, &isActive); err != nil {
			return nil, err
		}
		members = append(members, model.TeamMember{
			UserID:   uid,
			Username: username,
			IsActive: isActive,
		})
	}

	return &model.Team{
		TeamName: tn,
		Members:  members,
	}, nil
}

func (r *PostgresTeamRepo) Upsert(team *model.Team) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	_, err = tx.Exec(ctx, `INSERT INTO teams (team_name) VALUES ($1) ON CONFLICT (team_name) DO NOTHING`, team.TeamName)
	if err != nil {
		return err
	}

	for _, m := range team.Members {
		_, err := tx.Exec(ctx, `
			INSERT INTO users (user_id, username, team_name, skill_level, is_active)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (user_id) DO UPDATE
			  SET username = EXCLUDED.username,
			      team_name = EXCLUDED.team_name,
			      is_active = EXCLUDED.is_active
		`, m.UserID, m.Username, team.TeamName, model.SkillJunior.String(), m.IsActive)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (r *PostgresTeamRepo) List() ([]*model.Team, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.pool.Query(ctx, `SELECT team_name FROM teams`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*model.Team
	for rows.Next() {
		var tn string
		if err := rows.Scan(&tn); err != nil {
			return nil, err
		}
		t, err := r.GetByName(tn)
		if err != nil {
			// TODO
			continue
		}
		out = append(out, t)
	}
	return out, nil
}
