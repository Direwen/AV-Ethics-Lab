package response

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetCountBySessionID(ctx context.Context, sessionID string) (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetCountBySessionID(ctx context.Context, sessionID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Response{}).
		Joins("JOIN scenarios ON scenarios.id = responses.scenario_id").
		Where("scenarios.session_id = ?", sessionID).
		Count(&count).Error

	return count, err
}
