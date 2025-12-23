package scenario

import (
	"context"

	"github.com/direwen/go-server/internal/shared/domain"
	"github.com/direwen/go-server/internal/template"
)

type Service interface {
	GetNextScenario(ctx context.Context, sessionID string) (*Scenario, error)
}

type service struct {
	repo            Repository
	templateService template.Service
	llmClient       domain.LLMClient
}

func NewService(repo Repository, templateService template.Service, llmClient domain.LLMClient) Service {
	return &service{
		repo:            repo,
		templateService: templateService,
		llmClient:       llmClient,
	}
}

func (s *service) GetNextScenario(ctx context.Context, sessionID string) (*Scenario, error) {

	// Validate Session
	// Check Progress
	// Get used scenario context template ids
	// Get the unused context template
	// Generate the Scenario
	// Serialize Data
	// Save to DB

	return nil, nil
}
