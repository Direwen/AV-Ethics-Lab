package session

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/direwen/go-server/internal/shared/domain"
	"github.com/direwen/go-server/internal/util"
	"github.com/direwen/go-server/pkg/database"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Service interface {
	RegisterSession(ctx context.Context, input CreateSessionInput) (string, error)
	ValidateSession(ctx context.Context, session Session) error
	GetSession(ctx context.Context, sessionID uuid.UUID) (*Session, error)
	CompleteSession(ctx context.Context, session Session) error
	GetSessionFeedback(ctx context.Context, sessionID uuid.UUID) (*domain.FeedbackLLMResponse, error)
}

type service struct {
	repo    Repository
	llmPool domain.LLMPool
	experimentTargetCount int
}

func NewService(repo Repository, llmPool domain.LLMPool, experimentTargetCount int) Service {
	return &service{
		repo:    repo,
		llmPool: llmPool,
		experimentTargetCount: experimentTargetCount,
	}
}

func (s *service) RegisterSession(ctx context.Context, input CreateSessionInput) (string, error) {
	exists, err := s.repo.FingerprintExists(ctx, input.Fingerprint)
	if err != nil {
		return "", err
	}

	session_expiration_duration, err := time.ParseDuration(os.Getenv("SESSION_EXPIRATION"))
	if err != nil {
		return "", err
	}

	experimentPlan := domain.GenerateBalancedDesign(s.experimentTargetCount)
	planInJSON, err := json.Marshal(experimentPlan)
	if err != nil {
		return "", err
	}

	sess := Session{
		AgeRange:          input.AgeRange,
		Gender:            input.Gender,
		Country:           input.Country,
		Occupation:        input.Occupation,
		DrivingExperience: input.DrivingExperience,
		Fingerprint:       input.Fingerprint,
		SelfReportedNew:   input.SelfReportedNew,
		IsDuplicate:       exists,
		Status:            StatusActive,
		ExpiresAt:         time.Now().Add(session_expiration_duration),
		ExperimentPlan:    datatypes.JSON(planInJSON),
	}

	if err := s.repo.Create(ctx, &sess); err != nil {
		return "", err
	}

	signedToken, err := util.GenerateToken(
		sess.Id.String(),
		map[string]any{"issuer": "av-ethics-lab"},
	)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *service) ValidateSession(ctx context.Context, session Session) error {

	// Status Validation
	if session.Status != StatusActive {
		msg, found := SessionStatusErrorMsg[session.Status]
		if !found {
			msg = "session is not active"
		}
		return errors.New(msg)
	}

	// Expiry Validation
	if session.ExpiresAt.Before(time.Now()) {
		return errors.New("session expired")
	}

	return nil
}

func (s *service) GetSession(ctx context.Context, sessionID uuid.UUID) (*Session, error) {
	session, err := s.repo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return session, err
}

func (s *service) CompleteSession(ctx context.Context, session Session) error {
	session.Status = StatusCompleted
	return s.repo.Update(ctx, &session)
}

func (s *service) GetSessionFeedback(ctx context.Context, sessionID uuid.UUID) (*domain.FeedbackLLMResponse, error) {
	// Get session with all scenarios and responses preloaded
	session, err := s.repo.GetByID(ctx, sessionID, database.WithPreload("Scenarios.Response"))
	if err != nil {
		return nil, errors.New("session not found")
	}

	// Only completed sessions can get feedback
	if session.Status != StatusCompleted {
		return nil, errors.New("experiment not completed yet")
	}

	// Return cached feedback if exists
	if session.Feedback != nil {
		var feedback domain.FeedbackLLMResponse
		if err := json.Unmarshal(session.Feedback, &feedback); err != nil {
			return nil, errors.New("failed to parse cached feedback")
		}
		return &feedback, nil
	}

	// Build demographic data
	demographic := domain.Demographic{
		AgeRange:          session.AgeRange,
		Gender:            session.Gender,
		Country:           session.Country,
		Occupation:        session.Occupation,
		DrivingExperience: session.DrivingExperience,
	}

	// Build enriched responses from preloaded data
	responses := []domain.EnrichedResponse{}
	for _, scenario := range session.Scenarios {
		if scenario.Response == nil {
			continue
		}

		var factors domain.ScenarioFactors
		if err := json.Unmarshal(scenario.Factors, &factors); err != nil {
			continue
		}

		var rankedOptions []string
		if err := json.Unmarshal(scenario.Response.RankingOrder, &rankedOptions); err != nil {
			continue
		}

		responses = append(responses, domain.EnrichedResponse{
			Narrative:      scenario.Narrative,
			Factors:        factors,
			RankedOptions:  rankedOptions,
			HasInteracted:  scenario.Response.HasInteracted,
			ResponseTimeMs: scenario.Response.ResponseTimeMs,
			IsTimeout:      scenario.Response.IsTimeout,
		})
	}

	// Generate feedback via LLM
	result, err := s.llmPool.Execute(domain.TaskFeedback, func(client domain.Client) (any, error) {
		feedbackClient := client.(domain.FeedbackLLMClient)
		return feedbackClient.GenerateFeedback(ctx, domain.FeedbackLLMRequest{
			Demographic: demographic,
			Responses:   responses,
		})
	})
	if err != nil {
		return nil, err
	}
	feedback := result.(*domain.FeedbackLLMResponse)

	// Cache feedback in DB
	feedbackJSON, err := json.Marshal(feedback)
	if err != nil {
		return nil, err
	}
	session.Feedback = feedbackJSON
	if err := s.repo.Update(ctx, session); err != nil {
		return nil, err
	}

	return feedback, nil
}
