package usecase_tests

import (
	"testing"

	"github.com/leonidlivshits/Avito-MVP/internal/domain/model"
	"github.com/leonidlivshits/Avito-MVP/internal/domain/service"
	"github.com/leonidlivshits/Avito-MVP/internal/repository/inmemory"
	"github.com/leonidlivshits/Avito-MVP/internal/usecase"
)

func TestPRUsecase_CreateReassignMerge(t *testing.T) {
	userRepo := inmemory.NewInMemoryUserRepo()
	prRepo := inmemory.NewInMemoryPRRepo()
	teamRepo := inmemory.NewInMemoryTeamRepo()

	author := &model.User{UserID: "au", Username: "Author", TeamName: "teamA", Skill: model.SkillMiddle, IsActive: true}
	r1 := &model.User{UserID: "r1", Username: "R1", TeamName: "teamA", Skill: model.SkillMiddle, IsActive: true}
	r2 := &model.User{UserID: "r2", Username: "R2", TeamName: "teamA", Skill: model.SkillSenior, IsActive: true}
	r3 := &model.User{UserID: "r3", Username: "R3", TeamName: "teamA", Skill: model.SkillMiddle, IsActive: true}

	userRepo.Upsert(author)
	userRepo.Upsert(r1)
	userRepo.Upsert(r2)
	userRepo.Upsert(r3)

	assignSvc := service.NewAssignmentService(userRepo, prRepo)
	prUC := usecase.NewPRUsecase(prRepo, userRepo, teamRepo, assignSvc)

	// create PR
	pr, err := prUC.CreatePR("test-pr-1", "test pr", author.UserID)
	if err != nil {
		t.Fatalf("CreatePR failed: %v", err)
	}
	if pr.PullRequestID != "test-pr-1" {
		t.Fatalf("unexpected pr id")
	}

	// maybe assigned reviewers exist, try to reaasign
	if len(pr.AssignedReviewers) > 0 {
		old := pr.AssignedReviewers[0]
		newID, newPR, err := prUC.ReassignReviewer(pr.PullRequestID, old)
		if err != nil {
			t.Fatalf("Reassign failed: %v", err)
		}
		if newID == old {
			t.Fatalf("reassign didn't change reviewer")
		}
		_ = newPR
	}

	// merge
	merged, err := prUC.MergePR(pr.PullRequestID)
	if err != nil {
		t.Fatalf("MergePR failed: %v", err)
	}
	if merged.Status != model.PRStatusMerged {
		t.Fatalf("expected merged status")
	}
}
