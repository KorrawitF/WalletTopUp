package model

import (
	"WalletTopUp/internal/domain/entity"
	"time"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex"`
	FirstName    string
	LastName     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Wallet       *Wallet       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Transactions []Transaction `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (u *User) ToEntity() *entity.User {
	return &entity.User{
		ID:        u.ID,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}
