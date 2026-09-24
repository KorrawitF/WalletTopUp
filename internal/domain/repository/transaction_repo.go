package repository

import (
	"WalletTopUp/internal/domain/entity"
	"context"

	"gorm.io/gorm"
)

type TxRepo interface {
	WithTx(tx *gorm.DB) TxRepo
	Create(ctx context.Context, tx entity.Transaction) (*entity.Transaction, error)
	FindById(ctx context.Context, id string) (*entity.Transaction, error)
	Update(ctx context.Context, tx entity.Transaction) error
}
