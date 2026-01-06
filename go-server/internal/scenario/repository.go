package scenario

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, scenario *Scenario) error
	GetByContextTemplateID(ctx context.Context, id uuid.UUID) (*Scenario, error)
	GetUsedTemplateIDs(ctx context.Context, sessionID uuid.UUID) ([]uuid.UUID, error)
	GetPendingScenario(ctx context.Context, sessionID uuid.UUID) (*Scenario, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, scenario *Scenario) error {
	return r.db.WithContext(ctx).Create(scenario).Error
}

func (r *repository) GetByContextTemplateID(ctx context.Context, id uuid.UUID) (*Scenario, error) {
	var s Scenario

	err := r.db.WithContext(ctx).Preload("ContextTemplate").First(&s, "id = ?", id).Error
	return &s, err
}

func (r *repository) GetUsedTemplateIDs(ctx context.Context, sessionID uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID

	// Only count scenarios that have been answered (have a response)
	err := r.db.Model(&Scenario{}).
		Joins("INNER JOIN responses ON responses.scenario_id = scenarios.id").
		Where("scenarios.session_id = ?", sessionID).
		Pluck("scenarios.context_template_id", &ids).Error

	return ids, err
}

func (r *repository) GetPendingScenario(ctx context.Context, sessionID uuid.UUID) (*Scenario, error) {
	var s Scenario

	err := r.db.WithContext(ctx).
		Model(&Scenario{}).
		Joins("LEFT JOIN responses ON responses.scenario_id = scenarios.id").
		Where("scenarios.session_id = ?", sessionID).
		Where("responses.id IS NULL").
		Order("scenarios.created_at DESC").
		First(&s).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &s, nil
}
