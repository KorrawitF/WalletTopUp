package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTxManager_WithinTransaction_Commit(t *testing.T) {
	db, m := newMockDB(t, false)
	mgr := NewTxManager(db)

	m.ExpectBegin()
	m.ExpectCommit()

	called := false
	err := mgr.WithinTransaction(context.Background(), func(ctx context.Context) error {
		called = true
		assert.NotNil(t, mgr.GetTx(ctx), "tx should be stored in context")
		return nil
	})

	require.NoError(t, err)
	assert.True(t, called)
}

func TestTxManager_WithinTransaction_Rollback(t *testing.T) {
	db, m := newMockDB(t, false)
	mgr := NewTxManager(db)
	fnErr := errors.New("business failure")

	m.ExpectBegin()
	m.ExpectRollback()

	err := mgr.WithinTransaction(context.Background(), func(ctx context.Context) error {
		return fnErr
	})

	assert.ErrorIs(t, err, fnErr)
}

func TestTxManager_WithinTransaction_BeginFailed(t *testing.T) {
	db, m := newMockDB(t, false)
	mgr := NewTxManager(db)
	beginErr := errors.New("cannot begin")

	m.ExpectBegin().WillReturnError(beginErr)

	err := mgr.WithinTransaction(context.Background(), func(ctx context.Context) error {
		t.Fatal("fn must not run when begin fails")
		return nil
	})

	assert.ErrorIs(t, err, beginErr)
}

func TestTxManager_GetTx_NoTx(t *testing.T) {
	db, _ := newMockDB(t, false)

	assert.Nil(t, NewTxManager(db).GetTx(context.Background()))
}
