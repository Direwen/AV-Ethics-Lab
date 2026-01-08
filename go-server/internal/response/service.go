package response

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/direwen/go-server/internal/shared/services"
	"github.com/google/uuid"
)

type Service interface {
	SubmitResponse(ctx context.Context, sessionID, scenarioID uuid.UUID, input SubmitResponseInput) (*Response, error)
}

type service struct {
	repo            Repository
	sessionService  services.SessionValidator
	scenarioService services.ScenarioReader
}

func NewService(repo Repository, sessionService services.SessionValidator, scenarioService services.ScenarioReader) Service {
	return &service{
		repo:            repo,
		sessionService:  sessionService,
		scenarioService: scenarioService,
	}
}

func (s *service) SubmitResponse(ctx context.Context, sessionID, scenarioID uuid.UUID, input SubmitResponseInput) (*Response, error) {
	// Validate session is still active
	session, err := s.sessionService.GetSession(ctx, sessionID)
	if err != nil {
		return nil, errors.New("session not found")
	}

	if err := s.sessionService.ValidateSession(ctx, *session); err != nil {
		return nil, err
	}

	// Validate scenario belongs to session
	scenario, err := s.scenarioService.GetScenarioByID(ctx, scenarioID)
	if err != nil {
		return nil, errors.New("scenario not found")
	}

	if scenario.SessionID != sessionID {
		return nil, errors.New("scenario does not belong to this session")
	}

	// Check if response already exists
	existingResponse, err := s.repo.GetByScenarioID(ctx, scenarioID)
	if err == nil && existingResponse != nil {
		return nil, errors.New("response already exists")
	}

	// Save Response
	rankingOrderJSON, _ := json.Marshal(input.RankingOrder)
	response := &Response{
		ScenarioID:     scenarioID,
		IsTimeout:      input.IsTimeout,
		HasInteracted:  input.HasInteracted,
		ResponseTimeMs: input.ResponseTimeMs,
		RankingOrder:   rankingOrderJSON,
	}

	if err := s.repo.Create(ctx, response); err != nil {
		return nil, err
	}

	return response, nil
}
