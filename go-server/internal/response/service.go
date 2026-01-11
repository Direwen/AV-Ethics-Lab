package response

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/direwen/go-server/internal/shared/services"
	"github.com/direwen/go-server/pkg/database"
	"github.com/google/uuid"
)

type Service interface {
	SubmitResponse(ctx context.Context, sessionID, scenarioID uuid.UUID, input SubmitResponseInput) (*SubmitResponseOutput, error)
}

type service struct {
	repo               Repository
	sessionService     services.SessionValidator
	scenarioService    services.ScenarioReader
	transactionManager database.TransactionManager
}

func NewService(
	repo Repository,
	sessionService services.SessionValidator,
	scenarioService services.ScenarioReader,
	transactionManager database.TransactionManager,
) Service {
	return &service{
		repo:               repo,
		sessionService:     sessionService,
		scenarioService:    scenarioService,
		transactionManager: transactionManager,
	}
}

func (s *service) SubmitResponse(ctx context.Context, sessionID, scenarioID uuid.UUID, input SubmitResponseInput) (*SubmitResponseOutput, error) {
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
		maxAllowedTime := TimerDurationMs + NetworkBufferMs
		if actualElapsed > maxAllowedTime {
			input.IsTimeout = true
			input.ResponseTimeMs = TimerDurationMs
		}
	}

	// Check if response already exists (before transaction)
	existingResponse, err := s.repo.GetByScenarioID(ctx, scenarioID, database.WithSelect("id"))
	if err == nil && existingResponse != nil {
		return nil, errors.New("response already exists")
	}

	// Prepare response
	rankingOrderJSON, _ := json.Marshal(input.RankingOrder)
	response := &Response{
		ScenarioID:     scenarioID,
		IsTimeout:      input.IsTimeout,
		HasInteracted:  input.HasInteracted,
		ResponseTimeMs: input.ResponseTimeMs,
		RankingOrder:   rankingOrderJSON,
	}

	// Calculate completion before transaction
	var experimentPlan []any
	if err := json.Unmarshal(session.ExperimentPlan, &experimentPlan); err != nil {
		return nil, err
	}
	totalSteps := len(experimentPlan)

	var isComplete bool

	// Transaction: create response + count + update session status
	err = s.transactionManager.Do(ctx, func(txCtx context.Context) error {
		if err := s.repo.Create(txCtx, response); err != nil {
			return err
		}

		responseCount, err := s.repo.CountBySessionID(txCtx, sessionID)
		if err != nil {
			return err
		}

		isComplete = responseCount >= totalSteps
		if isComplete {
			if err := s.sessionService.CompleteSession(txCtx, *session); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &SubmitResponseOutput{
		Response:   response,
		IsComplete: isComplete,
	}, nil
}
