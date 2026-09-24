package model

import (
	"WalletTopUp/internal/domain/entity"
	"time"
)

type Wallet struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"not null;uniqueIndex"`
	Balance   float64 `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
	User      *User `gorm:"foreignKey:UserID"`
}

func (w *Wallet) ToEntity() *entity.Wallet {
	return &entity.Wallet{
		ID:      w.ID,
		UserID:  w.UserID,
		Balance: w.Balance,
	}
}
