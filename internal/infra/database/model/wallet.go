package model

import (
	"WalletTopUp/internal/domain/entity"
	"time"
)

type Wallet struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	Balance   float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (w *Wallet) ToEntity() *entity.Wallet {
	return &entity.Wallet{
		ID:      w.ID,
		UserID:  w.UserID,
		Balance: w.Balance,
	}
}
