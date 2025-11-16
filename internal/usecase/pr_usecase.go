package usecase

import (
	"time"

	"github.com/leonidlivshits/Avito-MVP/internal/domain"
	dservice "github.com/leonidlivshits/Avito-MVP/internal/domain/service"
	"github.com/leonidlivshits/Avito-MVP/internal/domain/model"
	"github.com/leonidlivshits/Avito-MVP/internal/ports"
)


type PRUsecase struct {
	prRepo        ports.PRRepository
	userRepo      ports.UserRepository
	teamRepo      ports.TeamRepository
	assignmentSvc *dservice.AssignmentService
}

func NewPRUsecase(pr ports.PRRepository, ur ports.UserRepository, tr ports.TeamRepository, as *dservice.AssignmentService) *PRUsecase {
	return &PRUsecase{
		prRepo:        pr,
		userRepo:      ur,
		teamRepo:      tr,
		assignmentSvc: as,
	}
}

func (u *PRUsecase) CreatePR(prID, prName, authorID string) (*model.PullRequest, error) {
	author, err := u.userRepo.GetByID(authorID)
	if err != nil {
		return nil, err
	}
	_ = author

	pr := &model.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   prName,
		AuthorID:          authorID,
		Status:            model.PRStatusOpen,
		AssignedReviewers: []string{},
		NeedMoreReviewers: false,
		CreatedAt:         time.Now().UTC(),
	}

	if _, err := u.assignmentSvc.AssignReviewsForPR(pr); err != nil {
		return nil, err
	}
	if err := u.prRepo.Create(pr); err != nil {
		return nil, err
	}
	return pr, nil
}

func (u *PRUsecase) ReassignReviewer(prID, oldUserID string) (string, *model.PullRequest, error) {
	pr, err := u.prRepo.GetByID(prID)
	if err != nil {
		return "", nil, err
	}
	if pr.Status == model.PRStatusMerged {
		return "", nil, domain.ErrPRMerged
	}
	found := false
	for _, r := range pr.AssignedReviewers {
		if r == oldUserID {
			found = true
			break
		}
	}
	if !found {
		return "", nil, domain.ErrNotAssigned
	}
	oldUser, err := u.userRepo.GetByID(oldUserID)
	if err != nil {
		return "", nil, err
	}
	author, err := u.userRepo.GetByID(pr.AuthorID)
	if err != nil {
		return "", nil, err
	}
	exclude := make(map[string]struct{})
	exclude[author.UserID] = struct{}{}
	for _, id := range pr.AssignedReviewers {
		exclude[id] = struct{}{}
	}
	candidates, err := u.userRepo.ListActiveByTeamExcluding(oldUser.TeamName, exclude)
	if err != nil {
		return "", nil, err
	}
	var chosen *model.User
	for _, c := range candidates {
		if c.Skill.AtLeast(author.Skill) {
			chosen = c
			break
		}
	}
	if chosen == nil {
		return "", nil, domain.ErrNoCandidate
	}
	for i, id := range pr.AssignedReviewers {
		if id == oldUserID {
			pr.AssignedReviewers[i] = chosen.UserID
			break
		}
	}
	if len(pr.AssignedReviewers) < 2 {
		remExcl := make(map[string]struct{})
		remExcl[author.UserID] = struct{}{}
		for _, id := range pr.AssignedReviewers {
			remExcl[id] = struct{}{}
		}
		remaining, _ := u.userRepo.ListActiveByTeamExcluding(oldUser.TeamName, remExcl)
		any := false
		for _, c := range remaining {
			if c.Skill.AtLeast(author.Skill) {
				any = true
				break
			}
		}
		pr.NeedMoreReviewers = any
	} else {
		pr.NeedMoreReviewers = false
	}
	if err := u.prRepo.Update(pr); err != nil {
		return "", nil, err
	}
	return chosen.UserID, pr, nil
}

func (u *PRUsecase) MergePR(prID string) (*model.PullRequest, error) {
	pr, err := u.prRepo.GetByID(prID)
	if err != nil {
		return nil, err
	}
	if pr.Status == model.PRStatusMerged {
		return pr, nil
	}
	now := time.Now().UTC()
	pr.Status = model.PRStatusMerged
	pr.MergedAt = &now
	if err := u.prRepo.Update(pr); err != nil {
		return nil, err
	}
	return pr, nil
}

func (u *PRUsecase) GetPR(prID string) (*model.PullRequest, error) {
	return u.prRepo.GetByID(prID)
}

func (u *PRUsecase) GetPRsForReviewer(userID string) ([]*model.PullRequest, error) {
	return u.prRepo.ListByReviewer(userID)
}

func (u *PRUsecase) ListOpen() ([]*model.PullRequest, error) {
	return u.prRepo.ListOpen()
}
