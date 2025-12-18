package session

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	FingerprintExists(ctx context.Context, fingerprint string) (bool, error)
	Create(ctx context.Context, session *Session) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FingerprintExists(ctx context.Context, fingerprint string) (bool, error) {
	var exists bool

	err := r.db.WithContext(ctx).Model(&Session{}).
		Select("1").
		Where("fingerprint = ?", fingerprint).
		Limit(1).
		Scan(&exists).Error

	return exists, err
}

func (r *repository) Create(ctx context.Context, session *Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}
