package model

type User struct {
	UserID   string     `json:"user_id"`
	Username string     `json:"username"`
	TeamName string     `json:"team_name"`
	Skill    SkillLevel `json:"skill_level"`
	IsActive bool       `json:"is_active"`
}
