package usecase

import (
	"github.com/leonidlivshits/Avito-MVP/internal/domain"
	"github.com/leonidlivshits/Avito-MVP/internal/domain/model"
	"github.com/leonidlivshits/Avito-MVP/internal/ports"
)

// управляет пользователями.
type UserUsecase struct {
	userRepo ports.UserRepository
	teamRepo ports.TeamRepository
}

func NewUserUsecase(ur ports.UserRepository, tr ports.TeamRepository) *UserUsecase {
	return &UserUsecase{userRepo: ur, teamRepo: tr}
}

// устанавливает флаг активности пользователя.
// Возвращает обновлённого пользователя.
func (u *UserUsecase) SetIsActive(userID string, isActive bool) (*model.User, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	user.IsActive = isActive
	if err := u.userRepo.Upsert(user); err != nil {
		return nil, err
	}
	return user, nil
}

// возвращает пользователя
func (u *UserUsecase) GetByID(userID string) (*model.User, error) {
	return u.userRepo.GetByID(userID)
}

// создаёт или обновляет пользователя с указанным уровнем навыка
// Проверяет, что команда существует.
func (u *UserUsecase) CreateOrUpdateUser(userID, username, teamName, skillText string, isActive bool) (*model.User, error) {
	if userID == "" || username == "" || teamName == "" || skillText == "" {
		return nil, domain.ErrInvalidRequest
	}

	if _, err := u.teamRepo.GetByName(teamName); err != nil {
		return nil, err
	}

	skill, perr := model.ParseSkillLevel(skillText)
	if perr != nil {
		return nil, domain.ErrInvalidRequest
	}

	user := &model.User{
		UserID:   userID,
		Username: username,
		TeamName: teamName,
		Skill:    skill,
		IsActive: isActive,
	}
	if err := u.userRepo.Upsert(user); err != nil {
		return nil, err
	}
	return user, nil
}

// возвращает список пользователей команды
func (u *UserUsecase) ListByTeam(teamName string) ([]*model.User, error) {
	return u.userRepo.ListByTeam(teamName)
}
