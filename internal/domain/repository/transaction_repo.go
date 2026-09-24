package repository

import (
	"WalletTopUp/internal/domain/entity"
	"context"
)

type TxRepo interface {
	Create(ctx context.Context, tx entity.Transaction) (*entity.Transaction, error)
}
