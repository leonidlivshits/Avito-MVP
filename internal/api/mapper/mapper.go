package mapper

import (
	"time"

	"github.com/leonidlivshits/Avito-MVP/internal/api/dto"
	"github.com/leonidlivshits/Avito-MVP/internal/domain/model"
)

func TeamDTOToDomain(in *dto.TeamDTO) *model.Team {
	if in == nil {
		return nil
	}
	out := &model.Team{
		TeamName: in.TeamName,
		Members:  make([]model.TeamMember, 0, len(in.Members)),
	}
	for _, m := range in.Members {
		out.Members = append(out.Members, model.TeamMember{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}
	return out
}

func TeamDomainToDTO(in *model.Team) *dto.TeamDTO {
	if in == nil {
		return nil
	}
	out := &dto.TeamDTO{
		TeamName: in.TeamName,
		Members:  make([]dto.TeamMemberDTO, 0, len(in.Members)),
	}
	for _, m := range in.Members {
		out.Members = append(out.Members, dto.TeamMemberDTO{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}
	return out
}

func UserDomainToDTO(in *model.User) *dto.UserDTO {
	if in == nil {
		return nil
	}
	return &dto.UserDTO{
		UserID:     in.UserID,
		Username:   in.Username,
		TeamName:   in.TeamName,
		SkillLevel: in.Skill.String(),
		IsActive:   in.IsActive,
	}
}

func PRDomainToDTO(in *model.PullRequest) *dto.PullRequestDTO {
	if in == nil {
		return nil
	}
	var mergedAt *string
	if in.MergedAt != nil {
		s := in.MergedAt.Format(time.RFC3339)
		mergedAt = &s
	}
	createdAt := in.CreatedAt.Format(time.RFC3339)
	return &dto.PullRequestDTO{
		PullRequestID:     in.PullRequestID,
		PullRequestName:   in.PullRequestName,
		AuthorID:          in.AuthorID,
		Status:            string(in.Status),
		AssignedReviewers: in.AssignedReviewers,
		NeedMoreReviewers: in.NeedMoreReviewers,
		CreatedAt:         createdAt,
		MergedAt:          mergedAt,
	}
}

func PRDomainToShortDTO(in *model.PullRequest) *dto.PullRequestShortDTO {
	if in == nil {
		return nil
	}
	return &dto.PullRequestShortDTO{
		PullRequestID:   in.PullRequestID,
		PullRequestName: in.PullRequestName,
		AuthorID:        in.AuthorID,
		Status:          string(in.Status),
	}
}
