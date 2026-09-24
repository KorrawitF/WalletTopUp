package database

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestMigrate_PropagatesError(t *testing.T) {
	sqlDB, m, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	migrateErr := errors.New("migration failed")
	m.MatchExpectationsInOrder(false)
	for range 10 {
		m.ExpectQuery(".*").WillReturnError(migrateErr)
		m.ExpectExec(".*").WillReturnError(migrateErr)
	}

	assert.Error(t, Migrate(db))
}
