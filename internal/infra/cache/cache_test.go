package cache

import (
	"WalletTopUp/internal/config"
	"WalletTopUp/internal/domain/entity"
	errs "WalletTopUp/internal/domain/error"
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRedis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = client.Close() })
	return mr, client
}

func sampleTx(expiresAt time.Time) *entity.Transaction {
	return &entity.Transaction{
		ID:            "tx-1",
		UserID:        1,
		Amount:        100,
		PaymentMethod: entity.PaymentMethodCreditCard,
		Status:        entity.TransactionStatusVerified,
		ExpiresAt:     expiresAt.UTC().Truncate(time.Second),
	}
}

func TestTxCache_SetAndGet(t *testing.T) {
	mr, client := newRedis(t)
	c := NewTxCache(client)
	ctx := context.Background()
	tx := sampleTx(time.Now().Add(10 * time.Minute))

	require.NoError(t, c.Set(ctx, tx))

	assert.True(t, mr.Exists(keyPrefix+"tx-1"))
	ttl := mr.TTL(keyPrefix + "tx-1")
	assert.True(t, ttl > 9*time.Minute && ttl <= 10*time.Minute, "unexpected ttl %s", ttl)

	got, err := c.Get(ctx, "tx-1")
	require.NoError(t, err)
	assert.Equal(t, tx, got)
}

func TestTxCache_Set_AlreadyExpiredIsSkipped(t *testing.T) {
	mr, client := newRedis(t)
	c := NewTxCache(client)

	require.NoError(t, c.Set(context.Background(), sampleTx(time.Now().Add(-time.Minute))))

	assert.False(t, mr.Exists(keyPrefix+"tx-1"))
}

func TestTxCache_Set_RedisError(t *testing.T) {
	mr, client := newRedis(t)
	c := NewTxCache(client)
	mr.SetError("boom")

	assert.Error(t, c.Set(context.Background(), sampleTx(time.Now().Add(time.Minute))))
}

func TestTxCache_Get_Miss(t *testing.T) {
	_, client := newRedis(t)
	c := NewTxCache(client)

	got, err := c.Get(context.Background(), "unknown")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, errs.ErrCacheMiss)
}

func TestTxCache_Get_RedisError(t *testing.T) {
	mr, client := newRedis(t)
	c := NewTxCache(client)
	mr.SetError("boom")

	got, err := c.Get(context.Background(), "tx-1")

	assert.Nil(t, got)
	assert.Error(t, err)
	assert.NotErrorIs(t, err, errs.ErrCacheMiss)
}

func TestTxCache_Get_InvalidPayload(t *testing.T) {
	mr, client := newRedis(t)
	c := NewTxCache(client)
	require.NoError(t, mr.Set(keyPrefix+"tx-1", "not-json"))

	got, err := c.Get(context.Background(), "tx-1")

	assert.Nil(t, got)
	assert.Error(t, err)
}

func TestTxCache_Delete(t *testing.T) {
	mr, client := newRedis(t)
	c := NewTxCache(client)
	require.NoError(t, mr.Set(keyPrefix+"tx-1", "{}"))

	require.NoError(t, c.Delete(context.Background(), "tx-1"))

	assert.False(t, mr.Exists(keyPrefix+"tx-1"))
}

func TestNoopTxCache(t *testing.T) {
	c := NewNoopTxCache()
	ctx := context.Background()

	assert.NoError(t, c.Set(ctx, sampleTx(time.Now())))
	got, err := c.Get(ctx, "tx-1")
	assert.Nil(t, got)
	assert.ErrorIs(t, err, errs.ErrCacheMiss)
	assert.NoError(t, c.Delete(ctx, "tx-1"))
}

func TestConnect_Success(t *testing.T) {
	mr := miniredis.RunT(t)

	client := Connect(context.Background(), config.CacheClient{Host: mr.Host(), Port: mr.Port()})

	require.NotNil(t, client)
	t.Cleanup(func() { _ = client.Close() })
	assert.NoError(t, client.Ping(context.Background()).Err())
}

func TestConnect_Failure(t *testing.T) {
	mr := miniredis.RunT(t)
	host, port := mr.Host(), mr.Port()
	mr.Close()

	assert.Nil(t, Connect(context.Background(), config.CacheClient{Host: host, Port: port}))
}
