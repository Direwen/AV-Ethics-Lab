package template

import (
	"context"
	"errors"
	"math/rand"
	"sync"

	"github.com/google/uuid"
)

type Service interface {
	LoadAllTemplates(ctx context.Context) error
	GetAllTemplates(ctx context.Context) ([]ContextTemplate, error)
	PickTemplate(excludeIDs []string) (*ContextTemplate, error)
}

type service struct {
	repo  Repository
	cache []ContextTemplate
	mu    sync.RWMutex //only one write | multiple reads
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) LoadAllTemplates(ctx context.Context) error {
	templates, err := s.repo.GetAll(ctx)

	if err != nil {
		return err
	}

	s.mu.Lock()
	s.cache = templates
	s.mu.Unlock()

	return nil
}

func (s *service) GetAllTemplates(ctx context.Context) ([]ContextTemplate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.cache) == 0 {
		return nil, errors.New("no templates found")
	}

	return s.cache, nil
}

func (s *service) PickTemplate(excludeIDs []string) (*ContextTemplate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.cache) == 0 {
		return nil, errors.New("template cache is empty")
	}

	var candidates []ContextTemplate
	excludeIDsMap := make(map[uuid.UUID]bool, len(excludeIDs))
	for _, id := range excludeIDs {
		parsedID, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		excludeIDsMap[parsedID] = true
	}

	for _, t := range s.cache {
		if _, found := excludeIDsMap[t.Id]; !found {
			candidates = append(candidates, t)
		}
	}

	if len(candidates) == 0 {
		candidates = s.cache
	}

	randomIndex := rand.Intn(len(candidates))
	return &candidates[randomIndex], nil
}
