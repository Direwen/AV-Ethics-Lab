package repository

import (
	"context"

	"github.com/direwen/go-server/internal/model"
	"gorm.io/gorm"
)

type TemplateRepository interface {
	Create(ctx context.Context, template *model.ContextTemplate) error
	FirstOrCreate(ctx context.Context, template *model.ContextTemplate) error
}

type templateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) TemplateRepository {
	return &templateRepository{db: db}
}

func (r *templateRepository) Create(ctx context.Context, template *model.ContextTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

func (r *templateRepository) FirstOrCreate(ctx context.Context, template *model.ContextTemplate) error {
	return r.db.WithContext(ctx).FirstOrCreate(template, model.ContextTemplate{Name: template.Name}).Error
}
