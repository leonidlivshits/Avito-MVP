package ports

import "github.com/leonidlivshits/Avito-MVP/internal/domain/model"

type TeamRepository interface {
	GetByName(teamName string) (*model.Team, error)
	Upsert(team *model.Team) error
	List() ([]*model.Team, error)
}
