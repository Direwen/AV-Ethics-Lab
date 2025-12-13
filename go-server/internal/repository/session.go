package repository

import (
	"context"

	"github.com/direwen/go-server/internal/model"
	"gorm.io/gorm"
)

type SessionRepository interface {
	FingerprintExists(ctx context.Context, fingerprint string) (bool, error)
	Create(ctx context.Context, session *model.Session) error
}

// Private Concrete Struct
type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	// Implicit Interface Compliance
	return &sessionRepository{db: db}
}

func (r *sessionRepository) FingerprintExists(ctx context.Context, fingerprint string) (bool, error) {
	var exists bool

	// Select("1") avoids fetching actual column data for the existence check
	err := r.db.WithContext(ctx).Model(&model.Session{}).
		Select("1").
		Where("fingerprint = ?", fingerprint).
		Limit(1).
		Scan(&exists).Error

	return exists, err
}

func (r *sessionRepository) Create(ctx context.Context, session *model.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}
