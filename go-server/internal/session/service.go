package session

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/direwen/go-server/internal/shared/domain"
	"github.com/direwen/go-server/internal/template"
	"github.com/direwen/go-server/internal/util"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Service interface {
	RegisterSession(ctx context.Context, input CreateSessionInput) (string, error)
	ValidateSession(ctx context.Context, session Session) error
	GetSession(ctx context.Context, id uuid.UUID) (*Session, error)
}

type service struct {
	repo            Repository
	templateService template.Service
}

func NewService(repo Repository, templateService template.Service) Service {
	return &service{
		repo:            repo,
		templateService: templateService,
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

	templates, err := s.templateService.GetAllTemplates(ctx)
	if err != nil {
		return "", err
	}
	experimentPlan := domain.GenerateBalancedDesign(len(templates))
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

func (s *service) GetSession(ctx context.Context, id uuid.UUID) (*Session, error) {
	session, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return session, err
}
