package repository

import (
	"WalletTopUp/internal/domain/entity"
	errs "WalletTopUp/internal/domain/error"
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const increaseBalanceSQL = `UPDATE "wallets" SET "balance"=balance + $1,"updated_at"=$2 WHERE user_id = $3 RETURNING *`

func TestWalletRepo_IncreaseBalance_Success(t *testing.T) {
	db, m := newMockDB(t, true)
	repo := NewWalletRepo(nil, db)

	m.ExpectQuery(q(increaseBalanceSQL)).
		WithArgs(50.0, sqlmock.AnyArg(), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}).AddRow(7, 1, 150.0))

	got, err := repo.IncreaseBalance(context.Background(), 1, 50)

	require.NoError(t, err)
	assert.Equal(t, &entity.Wallet{ID: 7, UserID: 1, Balance: 150}, got)
}

func TestWalletRepo_IncreaseBalance_NotFound(t *testing.T) {
	db, m := newMockDB(t, true)
	repo := NewWalletRepo(nil, db)

	m.ExpectQuery(q(increaseBalanceSQL)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance"}))

	got, err := repo.IncreaseBalance(context.Background(), 99, 50)

	assert.Nil(t, got)
	assert.ErrorIs(t, err, errs.ErrWalletNotFound)
}

func TestWalletRepo_IncreaseBalance_DBError(t *testing.T) {
	db, m := newMockDB(t, true)
	repo := NewWalletRepo(nil, db)
	dbErr := errors.New("deadlock")

	m.ExpectQuery(q(increaseBalanceSQL)).WillReturnError(dbErr)

	got, err := repo.IncreaseBalance(context.Background(), 1, 50)

	assert.Nil(t, got)
	assert.ErrorIs(t, err, dbErr)
}

func TestWalletRepo_WithTx(t *testing.T) {
	db, _ := newMockDB(t, true)
	other, _ := newMockDB(t, true)
	repo := NewWalletRepo(nil, db)

	assert.Same(t, other, repo.WithTx(other).(*walletRepo).db)
	assert.Same(t, db, repo.WithTx(nil).(*walletRepo).db, "nil tx should fall back to the base db")
}
