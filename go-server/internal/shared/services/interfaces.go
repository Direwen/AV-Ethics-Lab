package services

import (
	"context"

	"github.com/direwen/go-server/internal/shared/models"
	"github.com/google/uuid"
)

// SessionValidator provides session validation capabilities
type SessionValidator interface {
	GetSession(ctx context.Context, id uuid.UUID) (*models.Session, error)
	ValidateSession(ctx context.Context, session models.Session) error
	CompleteSession(ctx context.Context, session models.Session) error
}

// ScenarioReader provides read access to scenarios
type ScenarioReader interface {
	GetScenarioByID(ctx context.Context, id uuid.UUID) (*models.Scenario, error)
}
