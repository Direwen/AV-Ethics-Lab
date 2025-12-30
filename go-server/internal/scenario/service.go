package scenario

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/direwen/go-server/internal/session"
	"github.com/direwen/go-server/internal/shared/domain"
	"github.com/direwen/go-server/internal/template"
	"github.com/direwen/go-server/pkg/util"
	"github.com/google/uuid"
)

type Service interface {
	GetNextScenario(ctx context.Context, sessionID uuid.UUID) (*GetNextResponse, error)
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

func (s *service) GetNextScenario(ctx context.Context, sessionID uuid.UUID) (*GetNextResponse, error) {

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
	pendingScenario, err := s.repo.GetPendingScenario(ctx, sessionID)
	if err == nil && pendingScenario != nil {
		template, err := s.templateService.GetByID(pendingScenario.ContextTemplateID)
		if err != nil {
			return nil, err
		}
		return &GetNextResponse{
			Narrative: pendingScenario.Narrative,
			Entities:  pendingScenario.Entities,
			Factors:   pendingScenario.Factors,
			Width:     template.Width,
			Height:    template.Height,
			GridData:  template.GridData,
		}, nil
	}

	// Check Progress
	// Get used scenario context template ids
	usedContextIDs, err := s.repo.GetUsedTemplateIDs(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	currentStep := len(usedContextIDs)

	// Load the Experiment plan and keep the session status updated
	var experimentPlan []domain.ScenarioFactors
	if err := json.Unmarshal(session.ExperimentPlan, &experimentPlan); err != nil {
		return nil, errors.New("failed to load the experiment plan")
	}
	if currentStep >= len(experimentPlan) {
		// Mark Session Completed
		if err := s.sessionService.CompleteSession(ctx, *session); err != nil {
			return nil, err
		}
		return nil, errors.New("experiment completed")
	}

	// Pick a context template and factors for the current scenario
	contextTemplate, err := s.templateService.PickTemplate(usedContextIDs)
	if err != nil {
		return nil, err
	}
	currentFactors := experimentPlan[currentStep]

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

	// Enrich entities with IDs and Emojis
	enrichedEntities := make([]EnrichedEntity, len(llmRes.Entities))
	for i, e := range llmRes.Entities {
		info := domain.EntityRegistry[e.Type]
		enrichedEntities[i] = EnrichedEntity{
			ID:    fmt.Sprintf("ent_%s_%d", e.Type, i),
			Type:  e.Type,
			Emoji: info.Emoji,
			Row:   e.Row,
			Col:   e.Col,
			Metadata: EnrichedEntityMeta{
				IsStar:      e.Metadata.IsStar,
				IsViolation: e.Metadata.IsViolation,
				Action:      e.Metadata.Action,
			},
		}
	}

	// Serialize
	entitiesJSON, _ := json.Marshal(enrichedEntities)
	factorsJSON, _ := json.Marshal(currentFactors)
	newScenario := &Scenario{
		SessionID:         sessionID,
		Entities:          entitiesJSON,
		Factors:           factorsJSON,
		ContextTemplateID: contextTemplate.Id,
		Narrative:         llmRes.Narrative,
	}
	// Save to DB with retry
	if err := util.Retry(ctx, 3, 50*time.Millisecond, func() error {
		return s.repo.Create(ctx, newScenario)
	}); err != nil {
		return nil, err
	}
	res := &GetNextResponse{
		GridData:  contextTemplate.GridData,
		Entities:  entitiesJSON,
		Width:     contextTemplate.Width,
		Height:    contextTemplate.Height,
		Narrative: llmRes.Narrative,
		Factors:   factorsJSON,
	}

	return res, nil
}
