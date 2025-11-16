package usecase

import (
	"github.com/leonidlivshits/Avito-MVP/internal/domain"
	"github.com/leonidlivshits/Avito-MVP/internal/domain/model"
	"github.com/leonidlivshits/Avito-MVP/internal/ports"
)

type TeamUsecase struct {
	teamRepo ports.TeamRepository
	userRepo ports.UserRepository
}

func NewTeamUsecase(tr ports.TeamRepository, ur ports.UserRepository) *TeamUsecase {
	return &TeamUsecase{
		teamRepo: tr,
		userRepo: ur,
	}
}


func (u *TeamUsecase) CreateOrUpdateTeam(team *model.Team) error {
	if team == nil {
		return domain.ErrNotFound
	}

	if err := u.teamRepo.Upsert(team); err != nil {
		return err
	}

	for _, m := range team.Members {
		user := &model.User{
			UserID:   m.UserID,
			Username: m.Username,
			TeamName: team.TeamName,
			Skill:    model.SkillJunior,
			IsActive: m.IsActive,
		}
		if err := u.userRepo.Upsert(user); err != nil {
			return err
		}
	}

	return nil
}

func (u *TeamUsecase) GetTeam(teamName string) (*model.Team, error) {
	return u.teamRepo.GetByName(teamName)
}
