package inmemory

import (
	"sync"

	"github.com/leonidlivshits/Avito-MVP/internal/domain"
	"github.com/leonidlivshits/Avito-MVP/internal/domain/model"
	"github.com/leonidlivshits/Avito-MVP/internal/ports"
)

type InMemoryTeamRepo struct {
	mu    sync.RWMutex
	teams map[string]*model.Team
}

var _ ports.TeamRepository = (*InMemoryTeamRepo)(nil)

func NewInMemoryTeamRepo() *InMemoryTeamRepo {
	return &InMemoryTeamRepo{
		teams: make(map[string]*model.Team),
	}
}

func (r *InMemoryTeamRepo) GetByName(teamName string) (*model.Team, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.teams[teamName]
	if !ok {
		return nil, domain.ErrNotFound
	}
	c := *t
	return &c, nil
}

func (r *InMemoryTeamRepo) Upsert(team *model.Team) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c := *team
	r.teams[team.TeamName] = &c
	return nil
}

func (r *InMemoryTeamRepo) List() ([]*model.Team, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*model.Team, 0, len(r.teams))
	for _, t := range r.teams {
		c := *t
		out = append(out, &c)
	}
	return out, nil
}
