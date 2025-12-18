package scenario

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	GetLatestUnanswered(ctx context.Context, sessionID string) (*Scenario, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetLatestUnanswered(ctx context.Context, sessionID string) (*Scenario, error) {
	var s Scenario

	err := r.db.WithContext(ctx).
		Table("scenarios").
		Select("scenarios.*").
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
