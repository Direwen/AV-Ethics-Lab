package repository

import (
	"context"

	"github.com/direwen/go-server/internal/model"
	"gorm.io/gorm"
)

type ResponseRepository interface {
	GetCountBySessionID(ctx context.Context, sessionID string) (int64, error)
}

type responseRepository struct {
	db *gorm.DB
}

func NewResponseRepository(db *gorm.DB) ResponseRepository {
	return &responseRepository{db}
}

func (r *responseRepository) GetCountBySessionID(ctx context.Context, sessionID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Response{}).
		Joins("JOIN scenarios ON scenarios.id = responses.scenario_id").
		Where("scenarios.session_id = ?", sessionID).
		Count(&count).Error

	return count, err
}
