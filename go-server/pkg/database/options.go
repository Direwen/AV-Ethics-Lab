package database

import "gorm.io/gorm"

// Functional option for GORM queries
type QueryOption func(*gorm.DB) *gorm.DB

// To preload a relationship
func WithPreload(relation string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(relation)
	}
}

// To select specific fields
func WithSelect(fields ...string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(fields)
	}
}

// To add a custom WHERE clause
func WithFilter(condition string, args ...interface{}) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(condition, args...)
	}
}

// To apply all query options to a GORM DB instance
func ApplyOptions(db *gorm.DB, opts ...QueryOption) *gorm.DB {
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}
