package session

import (
	"context"

	"github.com/direwen/go-server/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	FingerprintExists(ctx context.Context, fingerprint string) (bool, error)
	Create(ctx context.Context, session *Session) error
	GetByID(ctx context.Context, id uuid.UUID, opts ...database.QueryOption) (*Session, error)
	Update(ctx context.Context, session *Session) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FingerprintExists(ctx context.Context, fingerprint string) (bool, error) {
	var exists bool

	err := database.GetDB(ctx, r.db).WithContext(ctx).Model(&Session{}).
		Select("1").
		Where("fingerprint = ?", fingerprint).
		Limit(1).
		Scan(&exists).Error

	return exists, err
}

func (r *repository) Create(ctx context.Context, session *Session) error {
	return database.GetDB(ctx, r.db).WithContext(ctx).Create(session).Error
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID, opts ...database.QueryOption) (*Session, error) {
	var session Session

	db := database.GetDB(ctx, r.db).WithContext(ctx).Model(&Session{}).Where("id = ?", id)
	db = database.ApplyOptions(db, opts...)

	err := db.First(&session).Error

	return &session, err
}

func (r *repository) Update(ctx context.Context, session *Session) error {
	return database.GetDB(ctx, r.db).WithContext(ctx).Save(session).Error
}
