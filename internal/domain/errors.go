package domain

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrPRExists       = errors.New("pr exists")
	ErrPRMerged       = errors.New("pr merged")
	ErrNotAssigned    = errors.New("not assigned")
	ErrNoCandidate    = errors.New("no candidate available")
	ErrInvalidRequest = errors.New("invalid request")
	ErrUnauthorized   = errors.New("unauthorized")
)
