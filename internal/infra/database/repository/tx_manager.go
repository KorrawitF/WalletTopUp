package repository

import (
	"WalletTopUp/internal/domain/repository"
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

type txManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) repository.TxManager {
	return &txManager{db: db}
}

func (m *txManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, txKey{}, tx))
	})
}

func (m *txManager) GetTx(ctx context.Context) *gorm.DB {
	val := ctx.Value(txKey{})

	tx, ok := val.(*gorm.DB)
	if !ok {
		return nil
	}

	return tx
}
