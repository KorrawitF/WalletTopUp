package service

import (
	"WalletTopUp/internal/domain/entity"
	"context"
)

type WalletSvc interface {
	VerifyTx(ctx context.Context, userId uint, amount float64, method string) (*entity.Transaction, error)
}
