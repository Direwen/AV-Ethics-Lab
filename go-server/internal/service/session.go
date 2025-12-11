package service

import (
	"context"
	"encoding/json"

	"github.com/direwen/go-server/internal/dto"
	"github.com/direwen/go-server/internal/model"
	"github.com/direwen/go-server/internal/repository"
	"github.com/direwen/go-server/internal/util"
	"gorm.io/datatypes"
)

type SessionService interface {
	RegisterSession(ctx context.Context, input dto.CreateSessionInput) (string, error)
}

// Private Concrete Struct
type sessionService struct {
	repo repository.SessionRepository
}

func NewSessionService(repo repository.SessionRepository) SessionService {
	// Implicit Interface Compliance
	return &sessionService{repo: repo}
}

func (s *sessionService) RegisterSession(ctx context.Context, input dto.CreateSessionInput) (string, error) {

	// Check if Fingerprint already exists
	exists, err := s.repo.FingerprintExists(ctx, input.Fingerprint)
	if err != nil {
		return "", err
	}

	// Encode demographics in JSON
	demographicsJSON, err := json.Marshal(input.Demographics)
	if err != nil {
		return "", err
	}

	session := model.Session{
		Demographics:    datatypes.JSON(demographicsJSON),
		Fingerprint:     input.Fingerprint,
		SelfReportedNew: input.SelfReportedNew,
		IsDuplicate:     exists,
	}

	// Create the session record, bind it to session
	if err := s.repo.Create(ctx, &session); err != nil {
		return "", err
	}

	// Generate JWT Token
	signedToken, err := util.GenerateToken(
		session.Id.String(),
		map[string]any{"issuer": "av-ethics-lab"},
	)
	if err != nil {
		return "", err
	}

	return signedToken, nil

}
