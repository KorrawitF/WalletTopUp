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

const (
	insertUserSQL = `INSERT INTO "users" ("username","first_name","last_name","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`
	selectUserSQL = `SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`
)

func TestUserRepo_Create(t *testing.T) {
	db, m := newMockDB(t, true)
	repo := NewUserRepo(nil, db)

	m.ExpectQuery(q(insertUserSQL)).
		WithArgs("john", "John", "Doe", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err := repo.Create(entity.User{Username: "john", FirstName: "John", LastName: "Doe"})

	assert.NoError(t, err)
}

func TestUserRepo_Create_Error(t *testing.T) {
	db, m := newMockDB(t, true)
	repo := NewUserRepo(nil, db)
	dbErr := errors.New("duplicate key")

	m.ExpectQuery(q(insertUserSQL)).WillReturnError(dbErr)

	assert.ErrorIs(t, repo.Create(entity.User{Username: "john"}), dbErr)
}

func TestUserRepo_FindById_Success(t *testing.T) {
	db, m := newMockDB(t, true)
	repo := NewUserRepo(nil, db)

	m.ExpectQuery(q(selectUserSQL)).WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "first_name", "last_name"}).
			AddRow(1, "john", "John", "Doe"))

	got, err := repo.FindById(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, &entity.User{ID: 1, Username: "john", FirstName: "John", LastName: "Doe"}, got)
}

func TestUserRepo_FindById_NotFound(t *testing.T) {
	db, m := newMockDB(t, true)
	repo := NewUserRepo(nil, db)

	m.ExpectQuery(q(selectUserSQL)).WillReturnRows(sqlmock.NewRows([]string{"id"}))

	got, err := repo.FindById(context.Background(), 99)

	assert.Nil(t, got)
	assert.ErrorIs(t, err, errs.ErrUserNotFound)
}

func TestUserRepo_FindById_DBError(t *testing.T) {
	db, m := newMockDB(t, true)
	repo := NewUserRepo(nil, db)
	dbErr := errors.New("connection reset")

	m.ExpectQuery(q(selectUserSQL)).WillReturnError(dbErr)

	got, err := repo.FindById(context.Background(), 1)

	assert.Nil(t, got)
	assert.ErrorIs(t, err, dbErr)
}
