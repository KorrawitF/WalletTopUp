package repository

import (
	"WalletTopUp/internal/domain/entity"
	errs "WalletTopUp/internal/domain/error"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	insertTxSQL = `INSERT INTO "transactions"`
	selectTxSQL = `SELECT * FROM "transactions" WHERE id = $1 ORDER BY "transactions"."id" LIMIT $2`
	updateTxSQL = `UPDATE "transactions" SET "status"=$1,"updated_at"=$2 WHERE "id" = $3`
)

func TestTxRepo_Create_Success(t *testing.T) {
	db, m := newMockDB(t, true)
	repo := NewTxRepo(nil, db)

	m.ExpectExec(q(insertTxSQL)).WillReturnResult(sqlmock.NewResult(0, 1))

	got, err := repo.Create(context.Background(), entity.Transaction{
		UserID:        1,
		Amount:        100,
		PaymentMethod: entity.PaymentMethodCreditCard,
		Status:        entity.TransactionStatusVerified,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, got.ID, "BeforeCreate should assign a UUID")
	assert.Equal(t, uint(1), got.UserID)
	assert.Equal(t, 100.0, got.Amount)
	assert.Equal(t, entity.PaymentMethodCreditCard, got.PaymentMethod)
	assert.Equal(t, entity.TransactionStatusVerified, got.Status)
	assert.WithinDuration(t, time.Now().Add(10*time.Minute), got.ExpiresAt, 5*time.Second)
}

func TestTxRepo_Create_Error(t *testing.T) {
	db, m := newMockDB(t, true)
	repo := NewTxRepo(nil, db)
	dbErr := errors.New("insert failed")

	m.ExpectExec(q(insertTxSQL)).WillReturnError(dbErr)

	got, err := repo.Create(context.Background(), entity.Transaction{UserID: 1})

	assert.Nil(t, got)
	assert.ErrorIs(t, err, dbErr)
}

func TestTxRepo_FindById_Success(t *testing.T) {
	db, m := newMockDB(t, true)
	repo := NewTxRepo(nil, db)
	expiresAt := time.Now().Add(5 * time.Minute).UTC()

	m.ExpectQuery(q(selectTxSQL)).WithArgs("tx-1", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "amount", "payment_method", "status", "expires_at"}).
			AddRow("tx-1", 1, 100.0, "debit_card", "verified", expiresAt))

	got, err := repo.FindById(context.Background(), "tx-1")

	require.NoError(t, err)
	assert.Equal(t, &entity.Transaction{
		ID:            "tx-1",
		UserID:        1,
		Amount:        100,
		PaymentMethod: entity.PaymentMethodDebitCard,
		Status:        entity.TransactionStatusVerified,
		ExpiresAt:     expiresAt,
	}, got)
}

func TestTxRepo_FindById_NotFound(t *testing.T) {
	db, m := newMockDB(t, true)
	repo := NewTxRepo(nil, db)

	m.ExpectQuery(q(selectTxSQL)).WithArgs("missing", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	got, err := repo.FindById(context.Background(), "missing")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, errs.ErrTxNotFound)
}

func TestTxRepo_FindById_DBError(t *testing.T) {
	db, m := newMockDB(t, true)
	repo := NewTxRepo(nil, db)
	dbErr := errors.New("connection reset")

	m.ExpectQuery(q(selectTxSQL)).WillReturnError(dbErr)

	got, err := repo.FindById(context.Background(), "tx-1")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, dbErr)
}

func TestTxRepo_Update(t *testing.T) {
	db, m := newMockDB(t, true)
	repo := NewTxRepo(nil, db)

	m.ExpectExec(q(updateTxSQL)).
		WithArgs("completed", sqlmock.AnyArg(), "tx-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(context.Background(), entity.Transaction{ID: "tx-1", Status: entity.TransactionStatusCompleted})

	assert.NoError(t, err)
}

func TestTxRepo_Update_Error(t *testing.T) {
	db, m := newMockDB(t, true)
	repo := NewTxRepo(nil, db)
	dbErr := errors.New("update failed")

	m.ExpectExec(q(updateTxSQL)).WillReturnError(dbErr)

	err := repo.Update(context.Background(), entity.Transaction{ID: "tx-1", Status: entity.TransactionStatusCompleted})

	assert.ErrorIs(t, err, dbErr)
}

func TestTxRepo_WithTx(t *testing.T) {
	db, _ := newMockDB(t, true)
	other, _ := newMockDB(t, true)
	repo := NewTxRepo(nil, db)

	assert.Same(t, other, repo.WithTx(other).(*txRepo).db)
	assert.Same(t, db, repo.WithTx(nil).(*txRepo).db, "nil tx should fall back to the base db")
}
