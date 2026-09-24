package handler

import (
	"WalletTopUp/internal/service"
	"WalletTopUp/pkg/lib"
)

type walletHandler struct {
	logger lib.Logger
	svc    service.TxSvc
}

func NewWalletHandler(logger lib.Logger, svc service.TxSvc) WalletHandler {
	return &walletHandler{
		logger,
		svc,
	}
}
