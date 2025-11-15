package inmemory

import (
	"sync"
	"time"

	"github.com/leonidlivshits/Avito-MVP/internal/domain"
	"github.com/leonidlivshits/Avito-MVP/internal/domain/model"
	"github.com/leonidlivshits/Avito-MVP/internal/ports"
)

type InMemoryPRRepo struct {
	mu sync.RWMutex
	pr map[string]*model.PullRequest
}

var _ ports.PRRepository = (*InMemoryPRRepo)(nil)

func NewInMemoryPRRepo() *InMemoryPRRepo {
	return &InMemoryPRRepo{
		pr: make(map[string]*model.PullRequest),
	}
}

func (r *InMemoryPRRepo) GetByID(prID string) (*model.PullRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.pr[prID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	c := *p
	return &c, nil
}

func (r *InMemoryPRRepo) Create(pr *model.PullRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.pr[pr.PullRequestID]; ok {
		return domain.ErrPRExists
	}
	p := *pr
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}
	r.pr[p.PullRequestID] = &p
	return nil
}

func (r *InMemoryPRRepo) Update(pr *model.PullRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.pr[pr.PullRequestID]; !ok {
		return domain.ErrNotFound
	}
	c := *pr
	r.pr[pr.PullRequestID] = &c
	return nil
}

func (r *InMemoryPRRepo) ListByReviewer(userID string) ([]*model.PullRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*model.PullRequest
	for _, p := range r.pr {
		for _, rev := range p.AssignedReviewers {
			if rev == userID {
				c := *p
				out = append(out, &c)
				break
			}
		}
	}
	return out, nil
}

func (r *InMemoryPRRepo) ListOpen() ([]*model.PullRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*model.PullRequest
	for _, p := range r.pr {
		if p.Status == model.PRStatusOpen {
			c := *p
			out = append(out, &c)
		}
	}
	return out, nil
}
