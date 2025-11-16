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


func RegisterHandlers(mux *http.ServeMux, logger *log.Logger, teamsUC *usecase.TeamUsecase, usersUC *usecase.UserUsecase, prsUC *usecase.PRUsecase, statsUC *usecase.StatsUsecase, adminToken, userToken string) {
	h := &Handlers{
		Teams: NewTeamsHandler(teamsUC),
		Users: NewUsersHandler(usersUC, prsUC, logger),
		PRs:   NewPRsHandler(prsUC, logger),
		Stats: NewStatsHandler(statsUC, logger),
	}

	mw := TokenAuthMiddleware(adminToken, userToken)

	mux.Handle("/team/add", mw(Adapt(h.Teams.CreateOrUpdateTeam)))
	mux.Handle("/team/get", mw(Adapt(h.Teams.GetTeam)))

	mux.Handle("/users/create", mw(Adapt(h.Users.CreateUser)))
	mux.Handle("/users/setIsActive", mw(Adapt(h.Users.SetIsActive)))
	mux.Handle("/users/getReview", mw(Adapt(h.Users.GetReview)))
	mux.Handle("/users/bulkDeactivate", mw(Adapt(h.Users.BulkDeactivate)))

	mux.Handle("/pullRequest/create", mw(Adapt(h.PRs.CreatePR)))
	mux.Handle("/pullRequest/get", mw(Adapt(h.PRs.GetPR)))
	mux.Handle("/pullRequest/merge", mw(Adapt(h.PRs.MergePR)))
	mux.Handle("/pullRequest/reassign", mw(Adapt(h.PRs.Reassign)))

	mux.Handle("/stats/assignments", mw(Adapt(h.Stats.Assignments)))
}
