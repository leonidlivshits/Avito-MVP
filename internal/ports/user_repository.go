package ports

import "github.com/leonidlivshits/Avito-MVP/internal/domain/model"

type UserRepository interface {
	GetByID(userID string) (*model.User, error)
	Upsert(user *model.User) error
	ListByTeam(teamName string) ([]*model.User, error)
	ListActiveByTeamExcluding(teamName string, excludeIDs map[string]struct{}) ([]*model.User, error)
	ListAll() ([]*model.User, error)
}
