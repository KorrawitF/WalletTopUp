package repository

import (
	"context"

	"gorm.io/gorm"
)

type TxManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	GetTx(ctx context.Context) *gorm.DB
}
