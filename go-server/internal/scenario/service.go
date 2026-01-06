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

	// Load the Experiment plan
	var experimentPlan []domain.ScenarioFactors
	if err := json.Unmarshal(session.ExperimentPlan, &experimentPlan); err != nil {
		return nil, errors.New("failed to load the experiment plan")
	}
	totalSteps := len(experimentPlan)

	// Get used scenario context template ids for progress tracking
	usedContextIDs, err := s.repo.GetUsedTemplateIDs(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Check the existence of the pending scenario
	pendingScenario, err := s.repo.GetPendingScenario(ctx, sessionID)
	if err == nil && pendingScenario != nil {
		tmpl, err := s.templateService.GetByID(pendingScenario.ContextTemplateID)
		if err != nil {
			return nil, err
		}
		var entities []EnrichedEntity
		if err := json.Unmarshal(pendingScenario.Entities, &entities); err != nil {
			return nil, err
		}
		var factors domain.ScenarioFactors
		if err := json.Unmarshal(pendingScenario.Factors, &factors); err != nil {
			return nil, err
		}
		var dilemmaOptions domain.DilemmaOptions
		if err := json.Unmarshal(pendingScenario.DilemmaOptions, &dilemmaOptions); err != nil {
			return nil, err
		}
		var gridData [][]int
		if err := json.Unmarshal(tmpl.GridData, &gridData); err != nil {
			return nil, err
		}
		var tridentSpawn domain.TridentSpawn
		if err := json.Unmarshal(pendingScenario.TridentSpawn, &tridentSpawn); err != nil {
			return nil, err
		}
		// Recalculate trident zones from stored spawn
		tridentZones := s.templateService.CalculateTridentZones(tmpl.Id, tridentSpawn)

		// Current step is the count of used templates + 1 (for the pending one)
		currentStep := len(usedContextIDs) + 1

		return &GetNextResponse{
			Narrative:      pendingScenario.Narrative,
			DilemmaOptions: dilemmaOptions,
			Entities:       entities,
			Factors:        factors,
			Width:          tmpl.Width,
			Height:         tmpl.Height,
			GridData:       gridData,
			LaneConfig:     s.templateService.GetLaneConfig(tmpl.Id),
			TridentZones:   tridentZones,
			TemplateName:   tmpl.Name,
			CurrentStep:    currentStep,
			TotalSteps:     totalSteps,
		}, nil
	}

	// Check Progress
	currentStep := len(usedContextIDs)

	if currentStep >= totalSteps {
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

	// Select a Trident Spawn point
	tridentSpawn, err := s.templateService.GetRandomTridentSpawn(contextTemplate.Id)
	if err != nil {
		return nil, errors.New("failed to get a trident spawn point")
	}
	// Calculate Trident Zones (with expandable B/C)
	tridentZones := s.templateService.CalculateTridentZones(contextTemplate.Id, *tridentSpawn)

	// Build Scenario LLM Request
	var gridData [][]int
	if err := json.Unmarshal(contextTemplate.GridData, &gridData); err != nil {
		return nil, err
	}
	laneConfig := s.templateService.GetLaneConfig(contextTemplate.Id)

	llmRes, err := s.llmClient.GenerateScenario(
		ctx,
		domain.ScenarioLLMRequest{
			TemplateName:   contextTemplate.Name,
			GridDimensions: fmt.Sprintf("%d:%d", contextTemplate.Width, contextTemplate.Height),
			Factors:        currentFactors,
			EgoPosition:    tridentSpawn.Coordinate,
			EgoOrientation: tridentSpawn.Orientation,
			TridentZones:   tridentZones,
		},
	)
	if err != nil {
		return nil, errors.New("failed to generate scenario")
	}

	// Add Ego AV entity (fixed position, not from LLM)
	egoEntity := EnrichedEntity{
		ID:    "ent_vehicle_av_ego",
		Type:  "vehicle_av",
		Emoji: domain.EntityRegistry["vehicle_av"].Emoji,
		Row:   tridentSpawn.Row,
		Col:   tridentSpawn.Col,
		Metadata: domain.EntityMeta{
			IsStar:      false,
			IsEgo:       true,
			IsViolation: false,
			Action:      "",
			Orientation: string(tridentSpawn.Orientation),
		},
	}

	// Enrich LLM entities with IDs and Emojis
	enrichedEntities := make([]EnrichedEntity, 0, len(llmRes.Entities)+1)
	enrichedEntities = append(enrichedEntities, egoEntity)
	for i, e := range llmRes.Entities {
		info := domain.EntityRegistry[e.Type]
		enrichedEntities = append(enrichedEntities, EnrichedEntity{
			ID:    fmt.Sprintf("ent_%s_%d", e.Type, i),
			Type:  e.Type,
			Emoji: info.Emoji,
			Row:   e.Row,
			Col:   e.Col,
			Metadata: domain.EntityMeta{
				IsStar:      e.Metadata.IsStar,
				IsEgo:       e.Metadata.IsEgo,
				IsViolation: e.Metadata.IsViolation,
				Action:      e.Metadata.Action,
				Orientation: e.Metadata.Orientation,
			},
		})
	}

	// Inject Tailgater if True
	if currentFactors.HasTailgater {
		rearCoord, err := s.templateService.GetRearCoordinate(contextTemplate.Id, tridentSpawn.Row, tridentSpawn.Col, tridentSpawn.Orientation)
		if err == nil {
			vehType := domain.CastRandomVehicle()
			tailgaterEntity := EnrichedEntity{
				ID:    "ent_" + vehType + "_tailgater",
				Type:  vehType,
				Emoji: domain.EntityRegistry[vehType].Emoji,
				Row:   rearCoord.Row,
				Col:   rearCoord.Col,
				Metadata: domain.EntityMeta{
					IsStar:      false,
					IsEgo:       false,
					IsViolation: true,
					Action:      "Tailgating dangerously close",
					Orientation: string(rearCoord.Orientation),
				},
			}
			enrichedEntities = append(enrichedEntities, tailgaterEntity)
		}
	}

	// Serialize for DB storage
	entitiesJSON, _ := json.Marshal(enrichedEntities)
	factorsJSON, _ := json.Marshal(currentFactors)
	dilemmaOptionsJSON, _ := json.Marshal(llmRes.DilemmaOptions)
	tridentSpawnJSON, _ := json.Marshal(tridentSpawn)
	newScenario := &Scenario{
		SessionID:         sessionID,
		Entities:          entitiesJSON,
		Factors:           factorsJSON,
		DilemmaOptions:    dilemmaOptionsJSON,
		ContextTemplateID: contextTemplate.Id,
		Narrative:         llmRes.Narrative,
		TridentSpawn:      tridentSpawnJSON,
	}
	// Save to DB with retry
	if err := util.Retry(ctx, 3, 50*time.Millisecond, func() error {
		return s.repo.Create(ctx, newScenario)
	}); err != nil {
		return nil, err
	}
	res := &GetNextResponse{
		GridData:       gridData,
		LaneConfig:     laneConfig,
		Entities:       enrichedEntities,
		Width:          contextTemplate.Width,
		Height:         contextTemplate.Height,
		Narrative:      llmRes.Narrative,
		DilemmaOptions: llmRes.DilemmaOptions,
		Factors:        currentFactors,
		TridentZones:   tridentZones,
		TemplateName:   contextTemplate.Name,
		CurrentStep:    currentStep + 1,
		TotalSteps:     totalSteps,
	}

	return res, nil
}
