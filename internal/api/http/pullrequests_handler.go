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

type PRsHandler struct {
	uc  *usecase.PRUsecase
	log *log.Logger
}

func NewPRsHandler(u *usecase.PRUsecase, logger *log.Logger) *PRsHandler {
	return &PRsHandler{uc: u, log: logger}
}

func (h *PRsHandler) CreatePR(w http.ResponseWriter, r *http.Request) error {
	var in dto.CreatePRRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return domain.ErrInvalidRequest
	}
	if in.PullRequestID == "" || in.AuthorID == "" {
		return domain.ErrInvalidRequest
	}
	pr, err := h.uc.CreatePR(in.PullRequestID, in.PullRequestName, in.AuthorID)
	if err != nil {
		return err
	}
	h.log.Infof("created pr %s by author %s", pr.PullRequestID, pr.AuthorID)
	return WriteJSON(w, http.StatusCreated, map[string]*dto.PullRequestDTO{"pr": mapper.PRDomainToDTO(pr)})
}

func (h *PRsHandler) GetPR(w http.ResponseWriter, r *http.Request) error {
	prID := r.URL.Query().Get("pull_request_id")
	if prID == "" {
		return domain.ErrInvalidRequest
	}
	pr, err := h.uc.GetPR(prID)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]*dto.PullRequestDTO{"pr": mapper.PRDomainToDTO(pr)})
}

func (h *PRsHandler) MergePR(w http.ResponseWriter, r *http.Request) error {
	var in dto.MergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return domain.ErrInvalidRequest
	}
	if in.PullRequestID == "" {
		return domain.ErrInvalidRequest
	}
	pr, err := h.uc.MergePR(in.PullRequestID)
	if err != nil {
		return err
	}
	h.log.Infof("merged pr %s", pr.PullRequestID)
	return WriteJSON(w, http.StatusOK, map[string]*dto.PullRequestDTO{"pr": mapper.PRDomainToDTO(pr)})
}

func (h *PRsHandler) Reassign(w http.ResponseWriter, r *http.Request) error {
	var in dto.ReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return domain.ErrInvalidRequest
	}
	if in.PullRequestID == "" || in.OldUserID == "" {
		return domain.ErrInvalidRequest
	}
	newID, pr, err := h.uc.ReassignReviewer(in.PullRequestID, in.OldUserID)
	if err != nil {
		return err
	}
	h.log.Infof("reassign pr=%s old=%s new=%s", pr.PullRequestID, in.OldUserID, newID)
	return WriteJSON(w, http.StatusOK, dto.ReassignResponse{
		Pr:         mapper.PRDomainToDTO(pr),
		ReplacedBy: newID,
	})
}
