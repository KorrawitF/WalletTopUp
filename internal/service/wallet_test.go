package service

import (
	cachemocks "WalletTopUp/internal/domain/cache/mocks"
	"WalletTopUp/internal/domain/entity"
	errs "WalletTopUp/internal/domain/error"
	repomocks "WalletTopUp/internal/domain/repository/mocks"
	libmocks "WalletTopUp/pkg/lib/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type svcDeps struct {
	logger     *libmocks.MockLogger
	userRepo   *repomocks.MockUserRepo
	txRepo     *repomocks.MockTxRepo
	txManager  *repomocks.MockTxManager
	walletRepo *repomocks.MockWalletRepo
	txCache    *cachemocks.MockTransactionCache
	svc        WalletSvc
}

func newSvcDeps(t *testing.T) *svcDeps {
	d := &svcDeps{
		logger:     libmocks.NewMockLogger(t),
		userRepo:   repomocks.NewMockUserRepo(t),
		txRepo:     repomocks.NewMockTxRepo(t),
		txManager:  repomocks.NewMockTxManager(t),
		walletRepo: repomocks.NewMockWalletRepo(t),
		txCache:    cachemocks.NewMockTransactionCache(t),
	}
	d.svc = NewWalletSvc(d.logger, d.userRepo, d.txRepo, d.txManager, d.walletRepo, d.txCache)
	return d
}

// runInTx makes the tx manager execute the callback and routes WithTx to the same repo mocks.
func (d *svcDeps) runInTx() {
	d.txManager.EXPECT().WithinTransaction(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		})
	d.txManager.EXPECT().GetTx(mock.Anything).Return(nil).Maybe()
	d.txRepo.EXPECT().WithTx(mock.Anything).Return(d.txRepo).Maybe()
	d.walletRepo.EXPECT().WithTx(mock.Anything).Return(d.walletRepo).Maybe()
}

func verifiedTx(expiresAt time.Time) *entity.Transaction {
	return &entity.Transaction{
		ID:            "tx-1",
		UserID:        1,
		Amount:        100,
		PaymentMethod: entity.PaymentMethodCreditCard,
		Status:        entity.TransactionStatusVerified,
		ExpiresAt:     expiresAt,
	}
}

func TestVerifyTx_Success(t *testing.T) {
	d := newSvcDeps(t)
	ctx := context.Background()
	created := verifiedTx(time.Now().Add(10 * time.Minute))

	d.userRepo.EXPECT().FindById(ctx, uint(1)).Return(&entity.User{ID: 1}, nil)
	d.txRepo.EXPECT().Create(ctx, entity.Transaction{
		UserID:        1,
		Amount:        100,
		PaymentMethod: entity.PaymentMethodCreditCard,
		Status:        entity.TransactionStatusVerified,
	}).Return(created, nil)
	d.txCache.EXPECT().Set(ctx, created).Return(nil)

	got, err := d.svc.VerifyTx(ctx, 1, 100, "credit_card")

	require.NoError(t, err)
	assert.Equal(t, created, got)
}

func TestVerifyTx_UserNotFound(t *testing.T) {
	d := newSvcDeps(t)
	ctx := context.Background()

	d.userRepo.EXPECT().FindById(ctx, uint(1)).Return(nil, errs.ErrUserNotFound)
	d.logger.EXPECT().Error(ctx, mock.Anything).Once()

	got, err := d.svc.VerifyTx(ctx, 1, 100, "credit_card")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, errs.ErrUserNotFound)
}

func TestVerifyTx_InvalidMethod(t *testing.T) {
	d := newSvcDeps(t)
	ctx := context.Background()

	d.userRepo.EXPECT().FindById(ctx, uint(1)).Return(&entity.User{ID: 1}, nil)
	d.logger.EXPECT().Error(ctx, mock.Anything).Once()

	got, err := d.svc.VerifyTx(ctx, 1, 100, "cash")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, errs.ErrInvalidMethod)
}

func TestVerifyTx_CreateFailed(t *testing.T) {
	d := newSvcDeps(t)
	ctx := context.Background()
	dbErr := errors.New("db down")

	d.userRepo.EXPECT().FindById(ctx, uint(1)).Return(&entity.User{ID: 1}, nil)
	d.txRepo.EXPECT().Create(ctx, mock.Anything).Return(nil, dbErr)
	d.logger.EXPECT().Error(ctx, mock.Anything).Once()

	got, err := d.svc.VerifyTx(ctx, 1, 100, "debit_card")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, dbErr)
}

func TestVerifyTx_CacheSetFailedStillSucceeds(t *testing.T) {
	d := newSvcDeps(t)
	ctx := context.Background()
	created := verifiedTx(time.Now().Add(10 * time.Minute))

	d.userRepo.EXPECT().FindById(ctx, uint(1)).Return(&entity.User{ID: 1}, nil)
	d.txRepo.EXPECT().Create(ctx, mock.Anything).Return(created, nil)
	d.txCache.EXPECT().Set(ctx, created).Return(errors.New("redis down"))
	d.logger.EXPECT().Warn(ctx, mock.Anything).Once()

	got, err := d.svc.VerifyTx(ctx, 1, 100, "bank_transfer")

	require.NoError(t, err)
	assert.Equal(t, created, got)
}

func TestConfirmTx_Success(t *testing.T) {
	d := newSvcDeps(t)
	ctx := context.Background()
	tx := verifiedTx(time.Now().Add(10 * time.Minute))

	d.txCache.EXPECT().Get(ctx, "tx-1").Return(nil, errs.ErrCacheMiss)
	d.runInTx()
	d.txRepo.EXPECT().FindById(mock.Anything, "tx-1").Return(tx, nil)
	d.walletRepo.EXPECT().IncreaseBalance(mock.Anything, uint(1), 100.0).Return(&entity.Wallet{UserID: 1, Balance: 250}, nil)
	d.txRepo.EXPECT().Update(mock.Anything, mock.MatchedBy(func(u entity.Transaction) bool {
		return u.Status == entity.TransactionStatusCompleted && u.CompletedAt != nil
	})).Return(nil)
	d.txCache.EXPECT().Delete(ctx, "tx-1").Return(nil)
	d.logger.EXPECT().Info(ctx, mock.Anything).Once()

	got, err := d.svc.ConfirmTx(ctx, "tx-1")

	require.NoError(t, err)
	assert.Equal(t, 250.0, got.Balance)
	assert.Equal(t, entity.TransactionStatusCompleted, got.Transaction.Status)
	assert.NotNil(t, got.Transaction.CompletedAt)
}

func TestConfirmTx_CachedTxExpired(t *testing.T) {
	d := newSvcDeps(t)
	ctx := context.Background()

	d.txCache.EXPECT().Get(ctx, "tx-1").Return(verifiedTx(time.Now().Add(-time.Minute)), nil)

	got, err := d.svc.ConfirmTx(ctx, "tx-1")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, errs.ErrTxExpired)
}

func TestConfirmTx_CachedTxValidProceeds(t *testing.T) {
	d := newSvcDeps(t)
	ctx := context.Background()
	tx := verifiedTx(time.Now().Add(10 * time.Minute))

	d.txCache.EXPECT().Get(ctx, "tx-1").Return(tx, nil)
	d.runInTx()
	d.txRepo.EXPECT().FindById(mock.Anything, "tx-1").Return(tx, nil)
	d.walletRepo.EXPECT().IncreaseBalance(mock.Anything, uint(1), 100.0).Return(&entity.Wallet{Balance: 100}, nil)
	d.txRepo.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)
	d.txCache.EXPECT().Delete(ctx, "tx-1").Return(nil)
	d.logger.EXPECT().Info(ctx, mock.Anything).Once()

	got, err := d.svc.ConfirmTx(ctx, "tx-1")

	require.NoError(t, err)
	assert.Equal(t, 100.0, got.Balance)
}

func TestConfirmTx_CacheReadErrorLogsAndProceeds(t *testing.T) {
	d := newSvcDeps(t)
	ctx := context.Background()
	tx := verifiedTx(time.Now().Add(10 * time.Minute))

	d.txCache.EXPECT().Get(ctx, "tx-1").Return(nil, errors.New("redis down"))
	d.runInTx()
	d.txRepo.EXPECT().FindById(mock.Anything, "tx-1").Return(tx, nil)
	d.walletRepo.EXPECT().IncreaseBalance(mock.Anything, uint(1), 100.0).Return(&entity.Wallet{Balance: 100}, nil)
	d.txRepo.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)
	d.txCache.EXPECT().Delete(ctx, "tx-1").Return(nil)
	d.logger.EXPECT().Warn(ctx, mock.Anything).Once()
	d.logger.EXPECT().Info(ctx, mock.Anything).Once()

	_, err := d.svc.ConfirmTx(ctx, "tx-1")

	require.NoError(t, err)
}

func TestConfirmTx_CacheDeleteErrorStillSucceeds(t *testing.T) {
	d := newSvcDeps(t)
	ctx := context.Background()
	tx := verifiedTx(time.Now().Add(10 * time.Minute))

	d.txCache.EXPECT().Get(ctx, "tx-1").Return(nil, errs.ErrCacheMiss)
	d.runInTx()
	d.txRepo.EXPECT().FindById(mock.Anything, "tx-1").Return(tx, nil)
	d.walletRepo.EXPECT().IncreaseBalance(mock.Anything, uint(1), 100.0).Return(&entity.Wallet{Balance: 100}, nil)
	d.txRepo.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)
	d.txCache.EXPECT().Delete(ctx, "tx-1").Return(errors.New("redis down"))
	d.logger.EXPECT().Warn(ctx, mock.Anything).Once()
	d.logger.EXPECT().Info(ctx, mock.Anything).Once()

	got, err := d.svc.ConfirmTx(ctx, "tx-1")

	require.NoError(t, err)
	assert.NotNil(t, got)
}

func TestConfirmTx_FindByIdFailed(t *testing.T) {
	d := newSvcDeps(t)
	ctx := context.Background()

	d.txCache.EXPECT().Get(ctx, "tx-1").Return(nil, errs.ErrCacheMiss)
	d.runInTx()
	d.txRepo.EXPECT().FindById(mock.Anything, "tx-1").Return(nil, errs.ErrTxNotFound)

	got, err := d.svc.ConfirmTx(ctx, "tx-1")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, errs.ErrTxNotFound)
}

func TestConfirmTx_InvalidStatus(t *testing.T) {
	cases := []struct {
		name   string
		status entity.TransactionStatus
		want   error
	}{
		{"completed", entity.TransactionStatusCompleted, errs.ErrTxCompleted},
		{"expired", entity.TransactionStatusExpired, errs.ErrTxExpired},
		{"unknown", entity.TransactionStatus("pending"), errs.ErrTxNotVerified},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := newSvcDeps(t)
			ctx := context.Background()
			tx := verifiedTx(time.Now().Add(10 * time.Minute))
			tx.Status = tc.status

			d.txCache.EXPECT().Get(ctx, "tx-1").Return(nil, errs.ErrCacheMiss)
			d.runInTx()
			d.txRepo.EXPECT().FindById(mock.Anything, "tx-1").Return(tx, nil)

			got, err := d.svc.ConfirmTx(ctx, "tx-1")

			assert.Nil(t, got)
			assert.ErrorIs(t, err, tc.want)
		})
	}
}

func TestConfirmTx_TxExpiredMarksExpired(t *testing.T) {
	d := newSvcDeps(t)
	ctx := context.Background()
	tx := verifiedTx(time.Now().Add(-time.Minute))

	d.txCache.EXPECT().Get(ctx, "tx-1").Return(nil, errs.ErrCacheMiss)
	d.runInTx()
	d.txRepo.EXPECT().FindById(mock.Anything, "tx-1").Return(tx, nil)
	d.txRepo.EXPECT().Update(mock.Anything, mock.MatchedBy(func(u entity.Transaction) bool {
		return u.Status == entity.TransactionStatusExpired
	})).Return(nil)
	d.txCache.EXPECT().Delete(ctx, "tx-1").Return(nil)

	got, err := d.svc.ConfirmTx(ctx, "tx-1")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, errs.ErrTxExpired)
}

func TestConfirmTx_TxExpiredUpdateFailed(t *testing.T) {
	d := newSvcDeps(t)
	ctx := context.Background()
	dbErr := errors.New("db down")

	d.txCache.EXPECT().Get(ctx, "tx-1").Return(nil, errs.ErrCacheMiss)
	d.runInTx()
	d.txRepo.EXPECT().FindById(mock.Anything, "tx-1").Return(verifiedTx(time.Now().Add(-time.Minute)), nil)
	d.txRepo.EXPECT().Update(mock.Anything, mock.Anything).Return(dbErr)
	d.logger.EXPECT().Error(mock.Anything, mock.Anything).Once()

	got, err := d.svc.ConfirmTx(ctx, "tx-1")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, dbErr)
}

func TestConfirmTx_IncreaseBalanceFailed(t *testing.T) {
	d := newSvcDeps(t)
	ctx := context.Background()

	d.txCache.EXPECT().Get(ctx, "tx-1").Return(nil, errs.ErrCacheMiss)
	d.runInTx()
	d.txRepo.EXPECT().FindById(mock.Anything, "tx-1").Return(verifiedTx(time.Now().Add(10*time.Minute)), nil)
	d.walletRepo.EXPECT().IncreaseBalance(mock.Anything, uint(1), 100.0).Return(nil, errs.ErrWalletNotFound)
	d.logger.EXPECT().Error(mock.Anything, mock.Anything).Once()

	got, err := d.svc.ConfirmTx(ctx, "tx-1")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, errs.ErrWalletNotFound)
}

func TestConfirmTx_UpdateCompletedFailed(t *testing.T) {
	d := newSvcDeps(t)
	ctx := context.Background()
	dbErr := errors.New("db down")

	d.txCache.EXPECT().Get(ctx, "tx-1").Return(nil, errs.ErrCacheMiss)
	d.runInTx()
	d.txRepo.EXPECT().FindById(mock.Anything, "tx-1").Return(verifiedTx(time.Now().Add(10*time.Minute)), nil)
	d.walletRepo.EXPECT().IncreaseBalance(mock.Anything, uint(1), 100.0).Return(&entity.Wallet{Balance: 100}, nil)
	d.txRepo.EXPECT().Update(mock.Anything, mock.Anything).Return(dbErr)
	d.logger.EXPECT().Error(mock.Anything, mock.Anything).Once()

	got, err := d.svc.ConfirmTx(ctx, "tx-1")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, dbErr)
}

func TestConfirmTx_BeginTransactionFailed(t *testing.T) {
	d := newSvcDeps(t)
	ctx := context.Background()
	dbErr := errors.New("cannot begin")

	d.txCache.EXPECT().Get(ctx, "tx-1").Return(nil, errs.ErrCacheMiss)
	d.txManager.EXPECT().WithinTransaction(ctx, mock.Anything).Return(dbErr)

	got, err := d.svc.ConfirmTx(ctx, "tx-1")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, dbErr)
}
