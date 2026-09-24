package mapper

import (
	"WalletTopUp/internal/domain/entity"
	"WalletTopUp/internal/infra/database/model"
)

func ToModel(entity entity.Wallet) *model.Wallet {
	return &model.Wallet{
		ID:      entity.ID,
		UserID:  entity.UserID,
		Balance: entity.Balance,
	}
}
