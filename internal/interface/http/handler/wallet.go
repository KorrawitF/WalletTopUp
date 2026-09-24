package handler

import (
	errs "WalletTopUp/internal/domain/error"
	"WalletTopUp/internal/interface/http/dto/request"
	"WalletTopUp/internal/interface/http/dto/response"
	"WalletTopUp/internal/service"
	"WalletTopUp/pkg/lib"
	"errors"
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
		var appErr *errs.Error
		if errors.As(err, &appErr) {
			c.JSON(appErr.Code, response.ErrorBody{
				Code:    appErr.Code,
				Status:  appErr.Status,
				Message: appErr.Message,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorBody{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "internal server error",
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

func (h *walletHandler) ConfirmTx(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.Confirm
	if err := c.BindJSON(&req); err != nil {
		h.logger.Error(ctx, lib.Meta{
			Event: "Confirm handler",
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

	res, err := h.svc.ConfirmTx(ctx, req.TransactionID)
	if err != nil {
		var appErr *errs.Error
		if errors.As(err, &appErr) {
			c.JSON(appErr.Code, response.ErrorBody{
				Code:    appErr.Code,
				Status:  appErr.Status,
				Message: appErr.Message,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorBody{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response.Confirm{
		TransactionID: res.Transaction.ID,
		UserID:        res.Transaction.UserID,
		Amount:        res.Transaction.Amount,
		Status:        string(res.Transaction.Status),
		Balance:       res.Balance,
	})
}
