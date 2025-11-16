package usecase

import "github.com/leonidlivshits/Avito-MVP/internal/ports"

type StatsUsecase struct {
	prRepo ports.PRRepository
}

func NewStatsUsecase(pr ports.PRRepository) *StatsUsecase {
	return &StatsUsecase{prRepo: pr}
}

func (s *StatsUsecase) CountAssignmentsPerUser() (map[string]int, error) {
	prs, err := s.prRepo.ListOpen()
	if err != nil {
		return nil, err
	}
	out := make(map[string]int)
	for _, p := range prs {
		for _, r := range p.AssignedReviewers {
			out[r]++
		}
	}
	return out, nil
}
