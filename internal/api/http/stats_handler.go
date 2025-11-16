package httpapi

import (
	"net/http"

	"github.com/leonidlivshits/Avito-MVP/internal/log"
	"github.com/leonidlivshits/Avito-MVP/internal/usecase"
)

type StatsHandler struct {
	uc  *usecase.StatsUsecase
	log *log.Logger
}

func NewStatsHandler(uc *usecase.StatsUsecase, logger *log.Logger) *StatsHandler {
	return &StatsHandler{uc: uc, log: logger}
}

func (h *StatsHandler) Assignments(w http.ResponseWriter, r *http.Request) error {
	m, err := h.uc.CountAssignmentsPerUser()
	if err != nil {
		return err
	}
	h.log.Infof("served stats assignments, users=%d", len(m))
	return WriteJSON(w, http.StatusOK, map[string]interface{}{"assignments": m})
}
