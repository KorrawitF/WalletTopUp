package model

import (
	"WalletTopUp/internal/domain/entity"
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"uniqueIndex"`
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) ToEntity() *entity.User {
	return &entity.User{
		ID:        u.ID,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}
