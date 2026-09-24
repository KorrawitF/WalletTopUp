package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// newMockDB returns a gorm DB (postgres dialect) backed by sqlmock.
func newMockDB(t *testing.T, skipDefaultTx bool) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, m, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, m.ExpectationsWereMet())
		_ = sqlDB.Close()
	})
	return openGorm(t, sqlDB, skipDefaultTx), m
}

func openGorm(t *testing.T, sqlDB *sql.DB, skipDefaultTx bool) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		SkipDefaultTransaction: skipDefaultTx,
		Logger:                 logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	return db
}

func q(sql string) string {
	return regexp.QuoteMeta(sql)
}
