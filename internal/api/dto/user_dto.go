package dto

type UserDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	SkillLevel string `json:"skill_level"`
	IsActive   bool   `json:"is_active"`
}

type CreateUserRequest struct {
	UserID     string `json:"user_id"`
	Username   string `json:"username"`
	TeamName   string `json:"team_name"`
	SkillLevel string `json:"skill_level"`
	IsActive   bool   `json:"is_active"`
}

type CreateUserResponse struct {
	User *UserDTO `json:"user"`
}

type SetIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type SetIsActiveResponse struct {
	User *UserDTO `json:"user"`
}

type BulkDeactivateRequest struct {
	TeamName string `json:"team_name"`
}
