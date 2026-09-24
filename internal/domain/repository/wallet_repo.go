package repository

import (
	"WalletTopUp/internal/domain/entity"
	"context"

	"gorm.io/gorm"
)

type WalletRepo interface {
	WithTx(tx *gorm.DB) WalletRepo
	IncreaseBalance(ctx context.Context, userId uint, amount float64) (*entity.Wallet, error)
}
