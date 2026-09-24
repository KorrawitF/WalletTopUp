package service

import (
	"WalletTopUp/internal/domain/entity"
	errs "WalletTopUp/internal/domain/error"
	"WalletTopUp/internal/domain/repository"
	"WalletTopUp/pkg/lib"
	"context"
)

type walletSvc struct {
	logger   lib.Logger
	userRepo repository.UserRepo
	txRepo   repository.TxRepo
}

func NewWalletSvc(logger lib.Logger, userRepo repository.UserRepo, txRepo repository.TxRepo) WalletSvc {
	return &walletSvc{
		logger,
		userRepo,
		txRepo,
	}
}

func (svc *walletSvc) VerifyTx(ctx context.Context, userId uint, amount float64, method string) (*entity.Transaction, error) {
	event := "Verify Tx service"

	user, err := svc.userRepo.FindById(ctx, userId)
	if err != nil {
		svc.logger.Error(ctx, lib.Meta{
			Event: event,
			Msg:   "failed to get user info.",
			Error: err,
		})
		return nil, err
	}

	if !entity.PaymentMethod(method).IsValid() {
		return nil, errs.ErrInvalidMethod
	}

	tx := entity.Transaction{
		UserID:        user.ID,
		Amount:        amount,
		PaymentMethod: entity.PaymentMethod(method),
		Status:        entity.TransactionStatusVerified,
	}

	createdTx, err := svc.txRepo.Create(ctx, tx)
	if err != nil {
		svc.logger.Error(ctx, lib.Meta{
			Event: event,
			Msg:   "failed during insert transaction.",
			Error: err,
		})
		return nil, err
	}

	return createdTx, nil
}
