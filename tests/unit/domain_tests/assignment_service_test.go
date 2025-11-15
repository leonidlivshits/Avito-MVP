package domain_tests

import (
	"testing"

	"github.com/leonidlivshits/Avito-MVP/internal/domain/model"
	"github.com/leonidlivshits/Avito-MVP/internal/domain/service"
	"github.com/leonidlivshits/Avito-MVP/internal/repository/inmemory"
)

func TestAssignReviewers_respectsSkillAndActiveAndNotAuthor(t *testing.T) {
	userRepo := inmemory.NewInMemoryUserRepo()
	prRepo := inmemory.NewInMemoryPRRepo()
	assigner := service.NewAssignmentService(userRepo, prRepo)

	author := &model.User{UserID: "u1", Username: "Author", TeamName: "backend", Skill: model.SkillMiddle, IsActive: true}
	u2 := &model.User{UserID: "u2", Username: "Low", TeamName: "backend", Skill: model.SkillJunior, IsActive: true}
	u3 := &model.User{UserID: "u3", Username: "Same", TeamName: "backend", Skill: model.SkillMiddle, IsActive: true}
	u4 := &model.User{UserID: "u4", Username: "High", TeamName: "backend", Skill: model.SkillSenior, IsActive: true}
	u5 := &model.User{UserID: "u5", Username: "InactiveHigh", TeamName: "backend", Skill: model.SkillSenior, IsActive: false}
	userRepo.Upsert(author)
	userRepo.Upsert(u2)
	userRepo.Upsert(u3)
	userRepo.Upsert(u4)
	userRepo.Upsert(u5)

	pr := &model.PullRequest{
		PullRequestID:    "pr1",
		PullRequestName:  "Add feature",
		AuthorID:         author.UserID,
		Status:           model.PRStatusOpen,
		AssignedReviewers: []string{},
	}

	selected, err := assigner.AssignReviewsForPR(pr)
	if err != nil {
		t.Fatalf("assign error: %v", err)
	}

	if len(selected) == 0 {
		t.Fatalf("expected at least one reviewer available (u3 or u4)")
	}
	for _, id := range selected {
		if id == "u1" || id == "u2" || id == "u5" {
			t.Fatalf("invalid selection included %s", id)
		}
	}
}
