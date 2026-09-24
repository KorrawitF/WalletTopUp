package repository

import (
	"WalletTopUp/internal/domain/entity"
	"WalletTopUp/internal/domain/repository"
	"WalletTopUp/internal/infra/database/mapper"
	"WalletTopUp/pkg/lib"
	"context"

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

func (repo *txRepo) Create(ctx context.Context, tx entity.Transaction) (*entity.Transaction, error) {
	model := mapper.ToTxModel(tx)
	if err := repo.db.Create(&model).Error; err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}
