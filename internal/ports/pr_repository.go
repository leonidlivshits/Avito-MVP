package ports

import "github.com/leonidlivshits/Avito-MVP/internal/domain/model"

type PRRepository interface {
	GetByID(prID string) (*model.PullRequest, error)
	Create(pr *model.PullRequest) error
	Update(pr *model.PullRequest) error
	ListByReviewer(userID string) ([]*model.PullRequest, error)
	ListOpen() ([]*model.PullRequest, error)
}