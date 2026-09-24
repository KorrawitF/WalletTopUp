package mapper

import (
	"WalletTopUp/internal/domain/entity"
	"WalletTopUp/internal/infra/database/model"
)

func ToTxModel(entity entity.Transaction) *model.Transaction {
	return &model.Transaction{
		ID:            entity.ID,
		UserID:        entity.UserID,
		Amount:        entity.Amount,
		PaymentMethod: string(entity.PaymentMethod),
		Status:        string(entity.Status),
	}
}
