package template

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, template *ContextTemplate) error
	FirstOrCreate(ctx context.Context, template *ContextTemplate) error
	GetAll(ctx context.Context) ([]ContextTemplate, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, template *ContextTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

func (r *repository) FirstOrCreate(ctx context.Context, template *ContextTemplate) error {
	return r.db.WithContext(ctx).FirstOrCreate(template, ContextTemplate{Name: template.Name}).Error
}

func (r *repository) GetAll(ctx context.Context) ([]ContextTemplate, error) {
	var templates []ContextTemplate
	err := r.db.WithContext(ctx).Find(&templates).Error
	return templates, err
}
