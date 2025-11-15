package service

import (
	"math/rand"
	"time"

	"github.com/leonidlivshits/Avito-MVP/internal/domain/model"
	"github.com/leonidlivshits/Avito-MVP/internal/ports"
	"github.com/leonidlivshits/Avito-MVP/internal/domain"
)

type AssignmentService struct {
	userRepo ports.UserRepository
	prRepo   ports.PRRepository
	seeded   *rand.Rand
}

func NewAssignmentService(u ports.UserRepository, p ports.PRRepository) *AssignmentService {
	return &AssignmentService{
		userRepo: u,
		prRepo:   p,
		seeded:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *AssignmentService) AssignReviewsForPR(pr *model.PullRequest) ([]string, error) {
	if pr == nil {
		return nil, domain.ErrNotFound
	}
	author, err := s.userRepo.GetByID(pr.AuthorID)
	if err != nil {
		return nil, err
	}

	if pr.Status == model.PRStatusMerged {
		return nil, domain.ErrPRMerged
	}

	exclude := make(map[string]struct{})
	exclude[author.UserID] = struct{}{}
	for _, a := range pr.AssignedReviewers {
		exclude[a] = struct{}{}
	}

	candidates, err := s.userRepo.ListActiveByTeamExcluding(author.TeamName, exclude)
	if err != nil {
		return nil, err
	}

	filtered := make([]*model.User, 0, len(candidates))
	for _, u := range candidates {
		if u.Skill.AtLeast(author.Skill) {
			filtered = append(filtered, u)
		}
	}

	s.seeded.Shuffle(len(filtered), func(i, j int) {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	})

	needed := 2 - len(pr.AssignedReviewers)
	if needed <= 0 {
		return pr.AssignedReviewers, nil
	}

	selected := make([]string, 0, needed)
	for i := 0; i < len(filtered) && len(selected) < needed; i++ {
		selected = append(selected, filtered[i].UserID)
	}

	pr.AssignedReviewers = append(pr.AssignedReviewers, selected...)
	if len(pr.AssignedReviewers) < 2 {
		remainingExcl := make(map[string]struct{})
		for _, id := range pr.AssignedReviewers {
			remainingExcl[id] = struct{}{}
		}
		remainingExcl[author.UserID] = struct{}{}
		remaining, _ := s.userRepo.ListActiveByTeamExcluding(author.TeamName, remainingExcl)
		anyRemaining := false
		for _, u := range remaining {
			if u.Skill.AtLeast(author.Skill) {
				anyRemaining = true
				break
			}
		}
		pr.NeedMoreReviewers = anyRemaining
	} else {
		pr.NeedMoreReviewers = false
	}

	return selected, nil
}
