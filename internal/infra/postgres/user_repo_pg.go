package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/leonidlivshits/Avito-MVP/internal/domain"
	"github.com/leonidlivshits/Avito-MVP/internal/domain/model"
	"github.com/leonidlivshits/Avito-MVP/internal/ports"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepo struct {
	pool *pgxpool.Pool
}

var _ ports.UserRepository = (*PostgresUserRepo)(nil)

func NewPostgresUserRepo(pool *pgxpool.Pool) *PostgresUserRepo {
	return &PostgresUserRepo{pool: pool}
}

func (r *PostgresUserRepo) GetByID(userID string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.pool.QueryRow(ctx, `
		SELECT user_id, username, team_name, skill_level, is_active
		FROM users WHERE user_id = $1
	`, userID)

	var id, username, teamName, skillText string
	var isActive bool
	if err := row.Scan(&id, &username, &teamName, &skillText, &isActive); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	skill, perr := model.ParseSkillLevel(skillText)
	if perr != nil {
		skill = model.SkillUnknown
	}
	return &model.User{
		UserID:   id,
		Username: username,
		TeamName: teamName,
		Skill:    skill,
		IsActive: isActive,
	}, nil
}

func (r *PostgresUserRepo) Upsert(user *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	skillText := user.Skill.String()

	_, err := r.pool.Exec(ctx, `
		INSERT INTO users (user_id, username, team_name, skill_level, is_active)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO UPDATE
		  SET username = EXCLUDED.username,
		      team_name = EXCLUDED.team_name,
		      skill_level = EXCLUDED.skill_level,
		      is_active = EXCLUDED.is_active
	`, user.UserID, user.Username, user.TeamName, skillText, user.IsActive)
	return err
}

func (r *PostgresUserRepo) ListByTeam(teamName string) ([]*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.pool.Query(ctx, `
		SELECT user_id, username, team_name, skill_level, is_active
		FROM users WHERE team_name = $1
	`, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*model.User
	for rows.Next() {
		var id, username, tn, skillText string
		var isActive bool
		if err := rows.Scan(&id, &username, &tn, &skillText, &isActive); err != nil {
			return nil, err
		}
		skill, perr := model.ParseSkillLevel(skillText)
		if perr != nil {
			skill = model.SkillUnknown
		}
		out = append(out, &model.User{
			UserID:   id,
			Username: username,
			TeamName: tn,
			Skill:    skill,
			IsActive: isActive,
		})
	}
	return out, nil
}

func (r *PostgresUserRepo) ListActiveByTeamExcluding(teamName string, excludeIDs map[string]struct{}) ([]*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []interface{}{teamName}
	query := `
		SELECT user_id, username, team_name, skill_level, is_active
		FROM users
		WHERE team_name = $1 AND is_active = true
	`
	if len(excludeIDs) > 0 {
		placeholders := make([]string, 0, len(excludeIDs))
		i := 2
		for id := range excludeIDs {
			args = append(args, id)
			placeholders = append(placeholders, fmt.Sprintf("$%d", i))
			i++
		}
		query += " AND user_id NOT IN (" + strings.Join(placeholders, ",") + ")"
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*model.User
	for rows.Next() {
		var id, username, tn, skillText string
		var isActive bool
		if err := rows.Scan(&id, &username, &tn, &skillText, &isActive); err != nil {
			return nil, err
		}
		skill, perr := model.ParseSkillLevel(skillText)
		if perr != nil {
			skill = model.SkillUnknown
		}
		out = append(out, &model.User{
			UserID:   id,
			Username: username,
			TeamName: tn,
			Skill:    skill,
			IsActive: isActive,
		})
	}
	return out, nil
}

func (r *PostgresUserRepo) ListAll() ([]*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.pool.Query(ctx, `
		SELECT user_id, username, team_name, skill_level, is_active
		FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*model.User
	for rows.Next() {
		var id, username, tn, skillText string
		var isActive bool
		if err := rows.Scan(&id, &username, &tn, &skillText, &isActive); err != nil {
			return nil, err
		}
		skill, perr := model.ParseSkillLevel(skillText)
		if perr != nil {
			skill = model.SkillUnknown
		}
		out = append(out, &model.User{
			UserID:   id,
			Username: username,
			TeamName: tn,
			Skill:    skill,
			IsActive: isActive,
		})
	}
	return out, nil
}
