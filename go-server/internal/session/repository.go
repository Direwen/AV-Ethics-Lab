package session

import (
	"context"
	"fmt"

	"github.com/direwen/go-server/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	FingerprintExists(ctx context.Context, fingerprint string) (bool, error)
	Create(ctx context.Context, session *Session) error
	GetByID(ctx context.Context, id uuid.UUID, opts ...database.QueryOption) (*Session, error)
	Update(ctx context.Context, session *Session) error
	CountSessions(ctx context.Context, opts ...database.QueryOption) (int64, error)
	CountCountries(ctx context.Context) (int64, error)
	CountArchetypes(ctx context.Context) ([]ArchetypeCount, error)
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

func (r *repository) CountSessions(ctx context.Context, opts ...database.QueryOption) (int64, error) {
	var count int64
	db := database.GetDB(ctx, r.db).WithContext(ctx).Model(&Session{})
	db = database.ApplyOptions(db, opts...)
	err := db.Count(&count).Error

	return count, err
}

func (r *repository) CountCountries(ctx context.Context) (int64, error) {
	var count int64
	err := database.GetDB(ctx, r.db).WithContext(ctx).
		Model(&Session{}).
		Distinct("country").
		Count(&count).Error

	return count, err
}

func (r *repository) CountArchetypes(ctx context.Context) ([]ArchetypeCount, error) {
	var counts []ArchetypeCount

	selector := "feedback->>'archetype'"

	err := database.GetDB(ctx, r.db).WithContext(ctx).
		Model(&Session{}).
		Select(fmt.Sprintf("%s as archetype, COUNT(*) as count", selector)).
		Group(selector).
		Order("count DESC").
		Scan(&counts).Error

	return counts, err
}
