package inmemory

import (
	"sync"

	"github.com/leonidlivshits/Avito-MVP/internal/domain"
	"github.com/leonidlivshits/Avito-MVP/internal/domain/model"
	"github.com/leonidlivshits/Avito-MVP/internal/ports"
)

type InMemoryUserRepo struct {
	mu    sync.RWMutex
	users map[string]*model.User
}

var _ ports.UserRepository = (*InMemoryUserRepo)(nil)

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		users: make(map[string]*model.User),
	}
}

func (r *InMemoryUserRepo) GetByID(userID string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[userID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	c := *u
	return &c, nil
}

func (r *InMemoryUserRepo) Upsert(user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c := *user
	r.users[user.UserID] = &c
	return nil
}

func (r *InMemoryUserRepo) ListByTeam(teamName string) ([]*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*model.User
	for _, u := range r.users {
		if u.TeamName == teamName {
			c := *u
			out = append(out, &c)
		}
	}
	return out, nil
}

func (r *InMemoryUserRepo) ListActiveByTeamExcluding(teamName string, excludeIDs map[string]struct{}) ([]*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*model.User
	for _, u := range r.users {
		if u.TeamName != teamName {
			continue
		}
		if !u.IsActive {
			continue
		}
		if _, excluded := excludeIDs[u.UserID]; excluded {
			continue
		}
		c := *u
		out = append(out, &c)
	}
	return out, nil
}

func (r *InMemoryUserRepo) ListAll() ([]*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*model.User, 0, len(r.users))
	for _, u := range r.users {
		c := *u
		out = append(out, &c)
	}
	return out, nil
}
