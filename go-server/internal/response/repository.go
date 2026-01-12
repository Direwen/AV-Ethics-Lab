package response

import (
	"context"

	"github.com/direwen/go-server/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, response *Response) error
	GetByID(ctx context.Context, id uuid.UUID, opts ...database.QueryOption) (*Response, error)
	GetByScenarioID(ctx context.Context, scenarioID uuid.UUID, opts ...database.QueryOption) (*Response, error)
	CountBySessionID(ctx context.Context, sessionID uuid.UUID) (int, error)
	GetBySessionID(ctx context.Context, sessionID uuid.UUID, opts ...database.QueryOption) ([]*Response, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, response *Response) error {
	return database.GetDB(ctx, r.db).WithContext(ctx).Create(response).Error
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID, opts ...database.QueryOption) (*Response, error) {
	var response Response
	db := database.GetDB(ctx, r.db).WithContext(ctx).Model(&Response{}).Where("id = ?", id)
	db = database.ApplyOptions(db, opts...)
	err := db.First(&response).Error
	return &response, err
}

func (r *repository) GetByScenarioID(ctx context.Context, scenarioID uuid.UUID, opts ...database.QueryOption) (*Response, error) {
	var response Response
	db := database.GetDB(ctx, r.db).WithContext(ctx).Model(&Response{}).Where("scenario_id = ?", scenarioID)
	db = database.ApplyOptions(db, opts...)
	err := db.First(&response).Error
	return &response, err
}

func (r *repository) CountBySessionID(ctx context.Context, sessionID uuid.UUID) (int, error) {
	var count int64
	err := database.GetDB(ctx, r.db).WithContext(ctx).
		Model(&Response{}).
		Joins("JOIN scenarios ON scenarios.id = responses.scenario_id").
		Where("scenarios.session_id = ?", sessionID).
		Count(&count).Error
	return int(count), err
}

func (r *repository) GetBySessionID(ctx context.Context, sessionID uuid.UUID, opts ...database.QueryOption) ([]*Response, error) {
	var responses []*Response
	db := database.GetDB(ctx, r.db).WithContext(ctx).
		Model(&Response{}).
		Joins("JOIN scenarios ON scenarios.id = responses.scenario_id").
		Where("scenarios.session_id = ?", sessionID)
	db = database.ApplyOptions(db, opts...)
	err := db.Find(&responses).Error
	return responses, err
}
