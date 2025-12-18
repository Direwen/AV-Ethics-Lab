package scenario

import (
	"context"
)

type Service interface {
	GetNextScenario(ctx context.Context, sessionID string) (*Scenario, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetNextScenario(ctx context.Context, sessionID string) (*Scenario, error) {
	scenario, err := s.repo.GetLatestUnanswered(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if scenario != nil {
		return scenario, nil
	}

	// If not, Generate a new scenario

	return nil, nil
}
