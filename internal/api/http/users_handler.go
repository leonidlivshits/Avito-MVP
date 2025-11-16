package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/leonidlivshits/Avito-MVP/internal/api/dto"
	"github.com/leonidlivshits/Avito-MVP/internal/api/mapper"
	"github.com/leonidlivshits/Avito-MVP/internal/domain"
	"github.com/leonidlivshits/Avito-MVP/internal/log"
	"github.com/leonidlivshits/Avito-MVP/internal/usecase"
)

type UsersHandler struct {
	uc   *usecase.UserUsecase
	prUC *usecase.PRUsecase
	log  *log.Logger
}

func NewUsersHandler(u *usecase.UserUsecase, prU *usecase.PRUsecase, logger *log.Logger) *UsersHandler {
	return &UsersHandler{uc: u, prUC: prU, log: logger}
}

func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	var in dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return domain.ErrInvalidRequest
	}
	if in.UserID == "" || in.Username == "" || in.TeamName == "" || in.SkillLevel == "" {
		return domain.ErrInvalidRequest
	}
	user, err := h.uc.CreateOrUpdateUser(in.UserID, in.Username, in.TeamName, in.SkillLevel, in.IsActive)
	if err != nil {
		return err
	}
	h.log.Infof("created user %s in team %s", user.UserID, user.TeamName)
	return WriteJSON(w, http.StatusCreated, map[string]*dto.UserDTO{"user": mapper.UserDomainToDTO(user)})
}

func (h *UsersHandler) SetIsActive(w http.ResponseWriter, r *http.Request) error {
	var in dto.SetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return domain.ErrInvalidRequest
	}
	if in.UserID == "" {
		return domain.ErrInvalidRequest
	}
	user, err := h.uc.SetIsActive(in.UserID, in.IsActive)
	if err != nil {
		return err
	}
	h.log.Infof("set isActive=%v for user %s", in.IsActive, in.UserID)
	return WriteJSON(w, http.StatusOK, map[string]*dto.UserDTO{"user": mapper.UserDomainToDTO(user)})
}

func (h *UsersHandler) GetReview(w http.ResponseWriter, r *http.Request) error {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		return domain.ErrInvalidRequest
	}
	prs, err := h.prUC.GetPRsForReviewer(userID)
	if err != nil {
		return err
	}
	shorts := make([]*dto.PullRequestShortDTO, 0, len(prs))
	for _, p := range prs {
		shorts = append(shorts, mapper.PRDomainToShortDTO(p))
	}
	resp := map[string]interface{}{
		"user_id":       userID,
		"pull_requests": shorts,
	}
	return WriteJSON(w, http.StatusOK, resp)
}
	

func (h *UsersHandler) BulkDeactivate(w http.ResponseWriter, r *http.Request) error {
	role := GetRole(r.Context())
	if role != domain.RoleAdmin {
		return domain.ErrUnauthorized
	}

	var in dto.BulkDeactivateRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return domain.ErrInvalidRequest
	}
	if in.TeamName == "" {
		return domain.ErrInvalidRequest
	}

	users, err := h.uc.ListByTeam(in.TeamName)
	if err != nil {
		return err
	}
	deactivated := make([]string, 0)
	for _, u := range users {
		if !u.IsActive {
			continue
		}
		if _, err := h.uc.SetIsActive(u.UserID, false); err != nil {
			h.log.Errorf("failed to deactivate user %s: %v", u.UserID, err)
			continue
		}
		deactivated = append(deactivated, u.UserID)
		h.log.Infof("deactivated user %s", u.UserID)
	}

	prs, err := h.prUC.ListOpen()
	if err != nil {
		h.log.Errorf("failed to list open prs: %v", err)
	}

	type reassignment struct {
		PRID    string `json:"pr_id"`
		Old     string `json:"old_user"`
		New     string `json:"new_user,omitempty"`
		Success bool   `json:"success"`
		Error   string `json:"error,omitempty"`
	}
	reassigned := make([]reassignment, 0)
	failed := make([]reassignment, 0)

	deactSet := make(map[string]struct{}, len(deactivated))
	for _, id := range deactivated {
		deactSet[id] = struct{}{}
	}

	for _, p := range prs {
		for _, ruid := range p.AssignedReviewers {
			if _, ok := deactSet[ruid]; !ok {
				continue
			}
			newID, _, err := h.prUC.ReassignReviewer(p.PullRequestID, ruid)
			if err != nil {
				failed = append(failed, reassignment{
					PRID:    p.PullRequestID,
					Old:     ruid,
					Success: false,
					Error:   err.Error(),
				})
				h.log.Warnf("failed reassign pr=%s old=%s err=%v", p.PullRequestID, ruid, err)
				continue
			}
			reassigned = append(reassigned, reassignment{
				PRID:    p.PullRequestID,
				Old:     ruid,
				New:     newID,
				Success: true,
			})
			h.log.Infof("reassigned pr=%s old=%s new=%s", p.PullRequestID, ruid, newID)
		}
	}

	resp := map[string]interface{}{
		"deactivated":     deactivated,
		"reassigned":      reassigned,
		"failed_reassign": failed,
	}
	return WriteJSON(w, http.StatusOK, resp)
}