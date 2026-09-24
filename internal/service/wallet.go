package service

import (
	"WalletTopUp/internal/domain/entity"
	errs "WalletTopUp/internal/domain/error"
	"WalletTopUp/internal/domain/repository"
	"WalletTopUp/internal/service/dto"
	"WalletTopUp/pkg/lib"
	"context"
	"fmt"
	"time"
)

type walletSvc struct {
	logger     lib.Logger
	userRepo   repository.UserRepo
	txRepo     repository.TxRepo
	txManager  repository.TxManager
	walletRepo repository.WalletRepo
}

func NewWalletSvc(logger lib.Logger, userRepo repository.UserRepo, txRepo repository.TxRepo, txManager repository.TxManager, walletRepo repository.WalletRepo) WalletSvc {
	return &walletSvc{
		logger,
		userRepo,
		txRepo,
		txManager,
		walletRepo,
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
		svc.logger.Error(ctx, lib.Meta{
			Event: event,
			Msg:   fmt.Sprintf("%s payment method not accepted.", method),
		})
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

func (svc *walletSvc) ConfirmTx(ctx context.Context, tx string) (*dto.ConfirmResult, error) {
	event := "Confirm Transaction"

	var (
		result     *dto.ConfirmResult
		expiredErr error
	)
	if err := svc.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		tTxRepo := svc.txRepo.WithTx(svc.txManager.GetTx(ctx))
		tx, err := tTxRepo.FindById(ctx, tx)
		if err != nil {
			return err
		}

		if err := validateStatus(*tx); err != nil {
			return err
		}

		now := time.Now()
		if tx.IsExpired(now) {
			tx.Status = entity.TransactionStatusExpired
			if err := tTxRepo.Update(ctx, *tx); err != nil {
				svc.logger.Error(ctx, lib.Meta{
					Event: event,
					Msg:   "failed to update expired tx",
					Error: err,
				})
				return err
			}

			expiredErr = errs.ErrTxExpired
			return nil
		}

		tWalletRepo := svc.walletRepo.WithTx(svc.txManager.GetTx(ctx))
		wallet, err := tWalletRepo.IncreaseBalance(ctx, tx.UserID, tx.Amount)
		if err != nil {
			svc.logger.Error(ctx, lib.Meta{
				Event: "Top up wallet balance",
				Msg:   "failed to top up wallet balance.",
				Error: err,
			})
			return err
		}

		tx.Status = entity.TransactionStatusCompleted
		tx.CompletedAt = &now
		if err := tTxRepo.Update(ctx, *tx); err != nil {
			svc.logger.Error(ctx, lib.Meta{
				Event: event,
				Msg:   "failed to update complated tx",
				Error: err,
			})
			return err
		}

		result = &dto.ConfirmResult{Transaction: tx, Balance: wallet.Balance}
		return nil
	}); err != nil {
		return nil, err
	}

	if expiredErr != nil {
		return nil, expiredErr
	}

	svc.logger.Info(ctx, lib.Meta{
		Event: event,
		Msg:   "Top up confirmed.",
	})

	return result, nil
}

func validateStatus(tx entity.Transaction) error {
	switch tx.Status {
	case entity.TransactionStatusCompleted:
		return errs.ErrTxCompleted
	case entity.TransactionStatusExpired:
		return errs.ErrTxExpired
	case entity.TransactionStatusVerified:
		return nil
	default:
		return errs.ErrTxNotVerified
	}
}
