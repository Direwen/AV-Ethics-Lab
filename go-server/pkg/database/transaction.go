package database

import (
	"context"

	"gorm.io/gorm"
)

// Private key type to prevent collisions with other context values
type TransactionKey struct{}

type TransactionManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type transactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &transactionManager{db: db}
}

// Wrap the function in a transaction and injects the tx (transaction) into the context
func (m *transactionManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {

	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Put the TX in the "Backpack" (Context)
		ctxWithTx := context.WithValue(ctx, TransactionKey{}, tx)
		// Run query with the new context
		return fn(ctxWithTx)
	})

}

// Helper to check the DB with transaction key in the context. If empty, return the fallback DB
func GetDB(ctx context.Context, fallbackDB *gorm.DB) *gorm.DB {
	if db, ok := ctx.Value(TransactionKey{}).(*gorm.DB); ok {
		return db
	}
	return fallbackDB
}
