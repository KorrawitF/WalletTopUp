package cache

import (
	"context"

	"WalletTopUp/internal/domain/cache"
	"WalletTopUp/internal/domain/entity"
	errs "WalletTopUp/internal/domain/error"
)

type NoopTransactionCache struct{}

func NewNoopTxCache() cache.TransactionCache {
	return &NoopTransactionCache{}
}

func (NoopTransactionCache) Set(context.Context, *entity.Transaction) error { return nil }

func (NoopTransactionCache) Get(context.Context, string) (*entity.Transaction, error) {
	return nil, errs.ErrCacheMiss
}

func (NoopTransactionCache) Delete(context.Context, string) error { return nil }
