package model

import "time"

type PRReviewer struct {
	PRID      string    `json:"pr_id"`
	UserID    string    `json:"user_id"`
	AssignedAt time.Time `json:"assigned_at"`
}
