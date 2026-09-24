package handler

import (
	"WalletTopUp/internal/domain/entity"
	errs "WalletTopUp/internal/domain/error"
	"WalletTopUp/internal/interface/http/dto/response"
	"WalletTopUp/internal/service/dto"
	svcmocks "WalletTopUp/internal/service/mocks"
	libmocks "WalletTopUp/pkg/lib/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupHandler(t *testing.T) (*gin.Engine, *svcmocks.MockWalletSvc, *libmocks.MockLogger) {
	logger := libmocks.NewMockLogger(t)
	svc := svcmocks.NewMockWalletSvc(t)
	h := NewWalletHandler(logger, svc)

	r := gin.New()
	r.POST("/verify", h.VerifyTx)
	r.POST("/confirm", h.ConfirmTx)
	return r, svc, logger
}

func doPost(r *gin.Engine, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func decode[T any](t *testing.T, w *httptest.ResponseRecorder) T {
	var v T
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &v))
	return v
}

func TestVerifyTx_Success(t *testing.T) {
	r, svc, _ := setupHandler(t)
	expiresAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)

	svc.EXPECT().VerifyTx(mock.Anything, uint(1), 100.0, "credit_card").Return(&entity.Transaction{
		ID:            "tx-1",
		UserID:        1,
		Amount:        100,
		PaymentMethod: entity.PaymentMethodCreditCard,
		Status:        entity.TransactionStatusVerified,
		ExpiresAt:     expiresAt,
	}, nil)

	w := doPost(r, "/verify", `{"user_id":1,"amount":100,"payment_method":"credit_card"}`)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, response.Verify{
		TransactionID: "tx-1",
		UserID:        1,
		Amount:        100,
		PaymentMethod: "credit_card",
		Status:        "verified",
		ExpiresAt:     expiresAt.Format(time.RFC3339),
	}, decode[response.Verify](t, w))
}

func TestVerifyTx_BadRequestBody(t *testing.T) {
	r, _, logger := setupHandler(t)
	logger.EXPECT().Error(mock.Anything, mock.Anything).Once()

	w := doPost(r, "/verify", `{invalid`)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	body := decode[response.ErrorBody](t, w)
	assert.Equal(t, "fail", body.Status)
	assert.Contains(t, body.Message, "failed to unmarshal request")
}

func TestVerifyTx_AppError(t *testing.T) {
	r, svc, _ := setupHandler(t)
	svc.EXPECT().VerifyTx(mock.Anything, uint(1), 100.0, "cash").Return(nil, errs.ErrInvalidMethod)

	w := doPost(r, "/verify", `{"user_id":1,"amount":100,"payment_method":"cash"}`)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, response.ErrorBody{
		Code:    http.StatusBadRequest,
		Status:  "fail",
		Message: errs.ErrInvalidMethod.Message,
	}, decode[response.ErrorBody](t, w))
}

func TestVerifyTx_InternalError(t *testing.T) {
	r, svc, _ := setupHandler(t)
	svc.EXPECT().VerifyTx(mock.Anything, uint(1), 100.0, "credit_card").Return(nil, errors.New("db down"))

	w := doPost(r, "/verify", `{"user_id":1,"amount":100,"payment_method":"credit_card"}`)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, response.ErrorBody{
		Code:    http.StatusInternalServerError,
		Status:  "error",
		Message: "internal server error",
	}, decode[response.ErrorBody](t, w))
}

func TestConfirmTx_Success(t *testing.T) {
	r, svc, _ := setupHandler(t)
	svc.EXPECT().ConfirmTx(mock.Anything, "tx-1").Return(&dto.ConfirmResult{
		Transaction: &entity.Transaction{
			ID:     "tx-1",
			UserID: 1,
			Amount: 100,
			Status: entity.TransactionStatusCompleted,
		},
		Balance: 300,
	}, nil)

	w := doPost(r, "/confirm", `{"transaction_id":"tx-1"}`)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, response.Confirm{
		TransactionID: "tx-1",
		UserID:        1,
		Amount:        100,
		Status:        "completed",
		Balance:       300,
	}, decode[response.Confirm](t, w))
}

func TestConfirmTx_BadRequestBody(t *testing.T) {
	r, _, logger := setupHandler(t)
	logger.EXPECT().Error(mock.Anything, mock.Anything).Once()

	w := doPost(r, "/confirm", `not-json`)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	body := decode[response.ErrorBody](t, w)
	assert.Equal(t, "fail", body.Status)
	assert.Contains(t, body.Message, "failed to unmarshal request")
}

func TestConfirmTx_AppError(t *testing.T) {
	r, svc, _ := setupHandler(t)
	svc.EXPECT().ConfirmTx(mock.Anything, "tx-1").Return(nil, errs.ErrTxExpired)

	w := doPost(r, "/confirm", `{"transaction_id":"tx-1"}`)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, errs.ErrTxExpired.Message, decode[response.ErrorBody](t, w).Message)
}

func TestConfirmTx_InternalError(t *testing.T) {
	r, svc, _ := setupHandler(t)
	svc.EXPECT().ConfirmTx(mock.Anything, "tx-1").Return(nil, errors.New("db down"))

	w := doPost(r, "/confirm", `{"transaction_id":"tx-1"}`)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "internal server error", decode[response.ErrorBody](t, w).Message)
}
