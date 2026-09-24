package dto

import (
	"WalletTopUp/internal/domain/entity"
)

type ConfirmResult struct {
	Transaction *entity.Transaction
	Balance     float64
}
