package scenario

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/direwen/go-server/internal/session"
	"github.com/direwen/go-server/internal/shared/domain"
	"github.com/direwen/go-server/internal/template"
	"github.com/google/uuid"
)

type Service interface {
	GetNextScenario(ctx context.Context, sessionID uuid.UUID) (*Scenario, error)
}

type service struct {
	repo            Repository
	sessionService  session.Service
	templateService template.Service
	llmClient       domain.LLMClient
}

func NewService(repo Repository, sessionService session.Service, templateService template.Service, llmClient domain.LLMClient) Service {
	return &service{
		repo:            repo,
		sessionService:  sessionService,
		templateService: templateService,
		llmClient:       llmClient,
	}
}

func (s *service) GetNextScenario(ctx context.Context, sessionID uuid.UUID) (*Scenario, error) {
	// Get Session
	session, err := s.sessionService.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Validate Session
	if err := s.sessionService.ValidateSession(ctx, *session); err != nil {
		return nil, err
	}

	// Check the existence of the pending scenario
	pending, err := s.repo.GetPendingScenario(ctx, sessionID)
	if err == nil && pending != nil {
		return pending, nil
	}

	// Generate the Scenario
	// Extract the experiment plan
	var experimentPlan []domain.ScenarioFactors
	if err := json.Unmarshal(session.ExperimentPlan, &experimentPlan); err != nil {
		return nil, errors.New("failed to load the experiment plan")
	}

	// Check Progress
	// Get used scenario context template ids
	usedContextIDs, err := s.repo.GetUsedTemplateIDs(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	currentStep := len(usedContextIDs)
	if currentStep >= len(experimentPlan) {
		return nil, errors.New("experiment completed")
	}
	currentFactors := experimentPlan[currentStep]

	// Get the unused context template
	contextTemplate, err := s.templateService.PickTemplate(usedContextIDs)
	if err != nil {
		return nil, err
	}

	// Build Scenario LLM Request struct
	var gridData [][]int
	if err := json.Unmarshal(contextTemplate.GridData, &gridData); err != nil {
		return nil, err
	}
	llmRes, err := s.llmClient.GenerateScenario(
		ctx,
		domain.ScenarioLLMRequest{
			TemplateName:   contextTemplate.Name,
			GridDimensions: fmt.Sprintf("%d:%d", contextTemplate.Width, contextTemplate.Height),
			GridData:       gridData,
			Factors:        currentFactors,
		},
	)
	if err != nil {
		return nil, errors.New("failed to generate scenario")
	}

	// Serialize
	entitiesJSON, _ := json.Marshal(llmRes.Entities)
	factorsJSON, _ := json.Marshal(currentFactors)
	newScenario := &Scenario{
		SessionID:         sessionID,
		Entities:          entitiesJSON,
		Factors:           factorsJSON,
		ContextTemplateID: contextTemplate.Id,
		Narrative:         llmRes.Narrative,
	}
	// Save to DB
	if err := s.repo.Create(ctx, newScenario); err != nil {
		return nil, err
	}

	return newScenario, nil
}
