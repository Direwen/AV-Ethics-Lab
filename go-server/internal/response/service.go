package response

import (
	"context"
	"encoding/json"
	"errors"
	"time"

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

	// Validate response time against scenario start time
	if scenario.StartedAt != nil {
		actualElapsed := time.Since(*scenario.StartedAt).Milliseconds()
		reportedTime := input.ResponseTimeMs

		// Maximum allowed time is timer duration + small buffer (15 seconds total)
		maxAllowedTime := int64(15000)

		// If actual elapsed time exceeds maximum, reject the response
		if actualElapsed > maxAllowedTime {
			return nil, errors.New("response time exceeded maximum allowed duration")
		}

		// Allow some tolerance for network latency (±2 seconds)
		tolerance := int64(2000)

		if reportedTime < actualElapsed-tolerance || reportedTime > actualElapsed+tolerance {
			// If time is suspicious, use the server-calculated time instead
			input.ResponseTimeMs = actualElapsed
		}

		// Ensure response time doesn't exceed maximum
		if input.ResponseTimeMs > maxAllowedTime {
			input.ResponseTimeMs = maxAllowedTime
		}
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
