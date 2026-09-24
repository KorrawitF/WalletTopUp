package repository

import (
	"WalletTopUp/internal/domain/entity"
	"context"
)

type UserRepo interface {
	Create(user entity.User) error
	FindById(ctx context.Context, id uint) (*entity.User, error)
}
