package repository

import (
	"WalletTopUp/internal/domain/entity"
	errs "WalletTopUp/internal/domain/error"
	"WalletTopUp/internal/domain/repository"
	"WalletTopUp/internal/infra/database/mapper"
	"WalletTopUp/internal/infra/database/model"
	"WalletTopUp/pkg/lib"
	"context"
	"errors"

	"gorm.io/gorm"
)

type txRepo struct {
	db *gorm.DB
}

func NewTxRepo(logger lib.Logger, db *gorm.DB) repository.TxRepo {
	return &txRepo{
		db,
	}
}

func (repo *txRepo) WithTx(tx *gorm.DB) repository.TxRepo {
	if tx == nil {
		tx = repo.db
	}
	return &txRepo{
		db: tx,
	}
}

func (repo *txRepo) Create(ctx context.Context, tx entity.Transaction) (*entity.Transaction, error) {
	model := mapper.ToTxModel(tx)
	if err := repo.db.WithContext(ctx).Create(&model).Error; err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

func (repo *txRepo) FindById(ctx context.Context, id string) (*entity.Transaction, error) {
	var model model.Transaction
	if err := repo.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errs.ErrTxNotFound
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

func (repo *txRepo) Update(ctx context.Context, tx entity.Transaction) error {
	model := mapper.ToTxModel(tx)
	return repo.db.WithContext(ctx).Updates(model).Error
}
