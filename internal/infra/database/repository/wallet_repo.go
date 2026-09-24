package repository

import (
	"WalletTopUp/internal/domain/entity"
	errs "WalletTopUp/internal/domain/error"
	"WalletTopUp/internal/domain/repository"
	"WalletTopUp/internal/infra/database/model"
	"WalletTopUp/pkg/lib"
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type walletRepo struct {
	db *gorm.DB
}

func NewWalletRepo(logger lib.Logger, db *gorm.DB) repository.WalletRepo {
	return &walletRepo{
		db,
	}
}

func (repo *walletRepo) WithTx(tx *gorm.DB) repository.WalletRepo {
	if tx == nil {
		tx = repo.db
	}

	return &walletRepo{
		db: tx,
	}
}

func (repo *walletRepo) IncreaseBalance(ctx context.Context, userId uint, amount float64) (*entity.Wallet, error) {
	var model model.Wallet
	res := repo.db.WithContext(ctx).Model(&model).
		Clauses(clause.Returning{}).
		Where("user_id = ?", userId).
		Update("balance", gorm.Expr("balance + ?", amount))
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, errs.ErrWalletNotFound
	}
	return model.ToEntity(), nil
}
