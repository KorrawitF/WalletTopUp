package service

import (
	"WalletTopUp/internal/domain/repository"
	"WalletTopUp/pkg/lib"
)

type txSvc struct {
	logger lib.Logger
	repo   repository.TxRepo
}

func NewTxSvc(logger lib.Logger, repo repository.TxRepo) TxSvc {
	return &txSvc{
		logger,
		repo,
	}
}
