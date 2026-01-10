package scenario

import (
	"context"
	"errors"

	"github.com/direwen/go-server/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, scenario *Scenario) error
	Update(ctx context.Context, scenario *Scenario) error
	GetByID(ctx context.Context, id uuid.UUID, opts ...database.QueryOption) (*Scenario, error)
	GetUsedTemplateIDs(ctx context.Context, sessionID uuid.UUID) ([]uuid.UUID, error)
	GetPendingScenario(ctx context.Context, sessionID uuid.UUID, opts ...database.QueryOption) (*Scenario, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, scenario *Scenario) error {
	return database.GetDB(ctx, r.db).WithContext(ctx).Create(scenario).Error
}

func (r *repository) Update(ctx context.Context, scenario *Scenario) error {
	return database.GetDB(ctx, r.db).WithContext(ctx).Save(scenario).Error
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID, opts ...database.QueryOption) (*Scenario, error) {
	var s Scenario
	db := database.GetDB(ctx, r.db).WithContext(ctx).Model(&Scenario{}).Where("id = ?", id)
	db = database.ApplyOptions(db, opts...)
	err := db.First(&s).Error
	return &s, err
}

func (r *repository) GetUsedTemplateIDs(ctx context.Context, sessionID uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID

	// Only count scenarios that have been answered (have a response)
	err := database.GetDB(ctx, r.db).Model(&Scenario{}).
		Joins("INNER JOIN responses ON responses.scenario_id = scenarios.id").
		Where("scenarios.session_id = ?", sessionID).
		Pluck("scenarios.context_template_id", &ids).Error

	return ids, err
}

func (r *repository) GetPendingScenario(ctx context.Context, sessionID uuid.UUID, opts ...database.QueryOption) (*Scenario, error) {
	var s Scenario

	db := database.GetDB(ctx, r.db).WithContext(ctx).
		Model(&Scenario{}).
		Joins("LEFT JOIN responses ON responses.scenario_id = scenarios.id").
		Where("scenarios.session_id = ?", sessionID).
		Where("responses.id IS NULL").
		Order("scenarios.created_at DESC")
	db = database.ApplyOptions(db, opts...)

	err := db.First(&s).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &s, nil
}
