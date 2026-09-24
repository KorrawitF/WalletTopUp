package handler

import (
	"WalletTopUp/internal/interface/http/dto/request"
	"WalletTopUp/internal/interface/http/dto/response"
	"WalletTopUp/internal/service"
	"WalletTopUp/pkg/lib"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type walletHandler struct {
	logger lib.Logger
	svc    service.WalletSvc
}

func NewWalletHandler(logger lib.Logger, svc service.WalletSvc) WalletHandler {
	return &walletHandler{
		logger,
		svc,
	}
}

func (h *walletHandler) VerifyTx(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.Verify
	if err := c.BindJSON(&req); err != nil {
		h.logger.Error(ctx, lib.Meta{
			Event: "Verify handler",
			Msg:   "failed to unmarshal request",
			Error: err,
		})
		c.JSON(http.StatusBadRequest, response.ErrorBody{
			Code:    http.StatusBadRequest,
			Status:  "fail",
			Message: fmt.Sprintf("failed to unmarshal request: %v", err),
		})
		return
	}

	tx, err := h.svc.VerifyTx(ctx, req.UserId, req.Amount, req.Method)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorBody{
			Code:    http.StatusBadRequest,
			Status:  "fail",
			Message: fmt.Sprintf("failed to verify tx: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, response.Verify{
		TransactionID: tx.ID,
		UserID:        tx.UserID,
		Amount:        tx.Amount,
		PaymentMethod: string(tx.PaymentMethod),
		Status:        string(tx.Status),
		ExpiresAt:     tx.ExpiresAt.Format(time.RFC3339),
	})
}
