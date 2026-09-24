package mapper

import (
	"WalletTopUp/internal/domain/entity"
	"WalletTopUp/internal/infra/database/model"
)

func ToUserModel(entity entity.User) *model.User {
	return &model.User{
		ID:        entity.ID,
		Username:  entity.Username,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
	}
}
