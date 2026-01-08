package response

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, response *Response) error
	GetByID(ctx context.Context, id uuid.UUID) (*Response, error)
	GetByScenarioID(ctx context.Context, scenarioID uuid.UUID) (*Response, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, response *Response) error {
	return r.db.WithContext(ctx).Create(response).Error
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Response, error) {
	var response Response
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&response).Error
	return &response, err
}

func (r *repository) GetByScenarioID(ctx context.Context, scenarioID uuid.UUID) (*Response, error) {
	var response Response
	err := r.db.WithContext(ctx).Where("scenario_id = ?", scenarioID).First(&response).Error
	return &response, err
}
