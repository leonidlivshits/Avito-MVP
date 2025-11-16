package httpapi

import (
	"net/http"

	"github.com/leonidlivshits/Avito-MVP/internal/api/dto"
	"github.com/leonidlivshits/Avito-MVP/internal/domain"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) error

func Adapt(h AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				_ = WriteError(w, http.StatusInternalServerError, dto.ErrorResponse{Error: dto.ErrorBody{Code: "INTERNAL", Message: "panic"}})
			}
		}()

		if err := h(w, r); err != nil {
			switch err {
			case domain.ErrNotFound:
				_ = WriteError(w, http.StatusNotFound, dto.ErrorResponse{Error: dto.ErrorBody{Code: "NOT_FOUND", Message: err.Error()}})
			case domain.ErrPRExists:
				_ = WriteError(w, http.StatusConflict, dto.ErrorResponse{Error: dto.ErrorBody{Code: "PR_EXISTS", Message: err.Error()}})
			case domain.ErrPRMerged:
				_ = WriteError(w, http.StatusConflict, dto.ErrorResponse{Error: dto.ErrorBody{Code: "PR_MERGED", Message: err.Error()}})
			case domain.ErrNotAssigned:
				_ = WriteError(w, http.StatusConflict, dto.ErrorResponse{Error: dto.ErrorBody{Code: "NOT_ASSIGNED", Message: err.Error()}})
			case domain.ErrNoCandidate:
				_ = WriteError(w, http.StatusConflict, dto.ErrorResponse{Error: dto.ErrorBody{Code: "NO_CANDIDATE", Message: err.Error()}})
			case domain.ErrInvalidRequest:
				_ = WriteError(w, http.StatusBadRequest, dto.ErrorResponse{Error: dto.ErrorBody{Code: "INVALID_REQUEST", Message: err.Error()}})
			case domain.ErrUnauthorized:
				_ = WriteError(w, http.StatusForbidden, dto.ErrorResponse{Error: dto.ErrorBody{Code: "UNAUTHORIZED", Message: err.Error()}})
			default:
				_ = WriteError(w, http.StatusInternalServerError, dto.ErrorResponse{Error: dto.ErrorBody{Code: "INTERNAL", Message: err.Error()}})
			}
		}
	}
}
