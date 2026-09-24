package repository

import (
	"WalletTopUp/internal/domain/repository"
	"WalletTopUp/pkg/lib"

	"gorm.io/gorm"
)

type txRepo struct {
	logger lib.Logger
	db     *gorm.DB
}

func NewTxRepo(logger lib.Logger, db *gorm.DB) repository.TxRepo {
	return &txRepo{
		logger,
		db,
	}
}
