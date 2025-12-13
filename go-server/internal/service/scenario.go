package service

import (
	"context"

	"github.com/direwen/go-server/internal/model"
	"github.com/direwen/go-server/internal/repository"
)

type ScenarioService interface {
	GetNextScenario(ctx context.Context, sessionID string) (*model.Scenario, error)
}

type scenarioService struct {
	repo repository.ScenarioRepository
}

func NewScenarioService(repo repository.ScenarioRepository) ScenarioService {
	return &scenarioService{repo: repo}
}

func (s *scenarioService) GetNextScenario(ctx context.Context, sessionID string) (*model.Scenario, error) {
	var scenario *model.Scenario

	// Get the latest unanswered scenario for the given session
	scenario, err := s.repo.GetLatestUnanswered(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// If exists, return the unanswered scenario
	if scenario != nil {
		return scenario, nil
	}

	// If not, Generate a new scenario

	return nil, nil
}
