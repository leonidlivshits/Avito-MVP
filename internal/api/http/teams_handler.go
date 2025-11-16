package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/leonidlivshits/Avito-MVP/internal/api/dto"
	"github.com/leonidlivshits/Avito-MVP/internal/api/mapper"
	"github.com/leonidlivshits/Avito-MVP/internal/usecase"
	"github.com/leonidlivshits/Avito-MVP/internal/domain"
)

type TeamsHandler struct {
	uc *usecase.TeamUsecase
}

func NewTeamsHandler(u *usecase.TeamUsecase) *TeamsHandler {
	return &TeamsHandler{uc: u}
}


func (h *TeamsHandler) CreateOrUpdateTeam(w http.ResponseWriter, r *http.Request) error {
	var in dto.TeamDTO
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return domain.ErrInvalidRequest
	}
	if in.TeamName == "" {
		return domain.ErrInvalidRequest
	}
	team := mapper.TeamDTOToDomain(&in)
	if team == nil {
		return domain.ErrInvalidRequest
	}
	if err := h.uc.CreateOrUpdateTeam(team); err != nil {
		return err
	}
	out := mapper.TeamDomainToDTO(team)
	return WriteJSON(w, http.StatusCreated, map[string]*dto.TeamDTO{"team": out})
}

func (h *TeamsHandler) GetTeam(w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query().Get("team_name")
	if q == "" {
		return domain.ErrInvalidRequest
	}
	team, err := h.uc.GetTeam(q)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, team)
}
