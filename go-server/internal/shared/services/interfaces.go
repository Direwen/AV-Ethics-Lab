package services

import (
	"context"

	"github.com/direwen/go-server/internal/scenario"
	"github.com/direwen/go-server/internal/session"
	"github.com/google/uuid"
)

// SessionValidator provides session validation capabilities
type SessionValidator interface {
	GetSession(ctx context.Context, id uuid.UUID) (*session.Session, error)
	ValidateSession(ctx context.Context, session session.Session) error
}

// ScenarioReader provides read access to scenarios
type ScenarioReader interface {
	GetScenarioByID(ctx context.Context, id uuid.UUID) (*scenario.Scenario, error)
}
