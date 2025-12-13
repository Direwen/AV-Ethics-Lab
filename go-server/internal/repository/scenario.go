package repository

import (
	"context"
	"errors"

	"github.com/direwen/go-server/internal/model"
	"gorm.io/gorm"
)

type ScenarioRepository interface {
	GetLatestUnanswered(ctx context.Context, sessionID string) (*model.Scenario, error)
}

type scenarioRepository struct {
	db *gorm.DB
}

func NewScenarioRepository(db *gorm.DB) ScenarioRepository {
	return &scenarioRepository{db}
}

func (r *scenarioRepository) GetLatestUnanswered(ctx context.Context, sessionID string) (*model.Scenario, error) {
	var scenario model.Scenario

	err := r.db.WithContext(ctx).
		Joins("LEFT JOIN responses ON scenarios.id = responses.scenario_id AND responses.session_id = ?", sessionID).
		Where("responses.id IS NULL").
		Order("scenarios.created_at DESC").
		First(&scenario).Error

	if err != nil {
		// No unanswered scenario, not an error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &scenario, nil
}
