package service

import (
	"WalletTopUp/internal/domain/entity"
	"WalletTopUp/internal/service/dto"
	"context"
)

type WalletSvc interface {
	VerifyTx(ctx context.Context, userId uint, amount float64, method string) (*entity.Transaction, error)
	ConfirmTx(ctx context.Context, tx string) (*dto.ConfirmResult, error)
}
