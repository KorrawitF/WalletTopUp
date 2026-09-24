package cache

import (
	"WalletTopUp/internal/domain/entity"
	"context"
)

type TransactionCache interface {
	Set(ctx context.Context, tx *entity.Transaction) error
	Get(ctx context.Context, id string) (*entity.Transaction, error)
	Delete(ctx context.Context, id string) error
}
