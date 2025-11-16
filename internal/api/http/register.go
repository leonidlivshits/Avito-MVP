package httpapi

import (
	"net/http"

	"github.com/leonidlivshits/Avito-MVP/internal/log"
	"github.com/leonidlivshits/Avito-MVP/internal/usecase"
)

type Handlers struct {
	Teams *TeamsHandler
	Users *UsersHandler
	PRs   *PRsHandler
	Stats *StatsHandler
}

func RegisterHandlers(mux *http.ServeMux, logger *log.Logger, teamsUC *usecase.TeamUsecase, usersUC *usecase.UserUsecase, prsUC *usecase.PRUsecase, statsUC *usecase.StatsUsecase, adminToken, userToken, authorToken, reviewerToken string) {
	h := &Handlers{
		Teams: NewTeamsHandler(teamsUC),
		Users: NewUsersHandler(usersUC, prsUC, logger),
		PRs:   NewPRsHandler(prsUC, logger),
		Stats: NewStatsHandler(statsUC, logger),
	}

	authMW := TokenAuthMiddleware(adminToken, userToken, authorToken, reviewerToken)

	mux.Handle("/team/add", authMW(Adapt(h.Teams.CreateOrUpdateTeam)))
	mux.Handle("/team/get", authMW(Adapt(h.Teams.GetTeam)))

	mux.Handle("/users/create", authMW(Adapt(h.Users.CreateUser)))
	mux.Handle("/users/setIsActive", authMW(Adapt(h.Users.SetIsActive)))
	mux.Handle("/users/getReview", authMW(Adapt(h.Users.GetReview)))
	mux.Handle("/users/bulkDeactivate", authMW(Adapt(h.Users.BulkDeactivate)))

	mux.Handle("/pullRequest/create", authMW(Adapt(h.PRs.CreatePR)))
	mux.Handle("/pullRequest/get", authMW(Adapt(h.PRs.GetPR)))
	mux.Handle("/pullRequest/merge", authMW(Adapt(h.PRs.MergePR)))
	mux.Handle("/pullRequest/reassign", authMW(Adapt(h.PRs.Reassign)))

	mux.Handle("/stats/assignments", authMW(Adapt(h.Stats.Assignments)))
}
