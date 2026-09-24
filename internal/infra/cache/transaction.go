package cache

import (
	"WalletTopUp/internal/domain/cache"
	"WalletTopUp/internal/domain/entity"
	errs "WalletTopUp/internal/domain/error"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const keyPrefix = "wallet:topup:tx:"

type transactionCache struct {
	client *redis.Client
}

func NewTxCache(client *redis.Client) cache.TransactionCache {
	return &transactionCache{
		client,
	}
}

func (c *transactionCache) Set(ctx context.Context, tx *entity.Transaction) error {
	ttl := time.Until(tx.ExpiresAt)
	if ttl <= 0 {
		return nil
	}
	data, err := json.Marshal(tx)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, keyPrefix+tx.ID, data, ttl).Err()
}

func (c *transactionCache) Get(ctx context.Context, id string) (*entity.Transaction, error) {
	data, err := c.client.Get(ctx, keyPrefix+id).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, errs.ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}
	var tx entity.Transaction
	if err := json.Unmarshal(data, &tx); err != nil {
		return nil, err
	}
	return &tx, nil
}

func (c *transactionCache) Delete(ctx context.Context, id string) error {
	return c.client.Del(ctx, keyPrefix+id).Err()
}
