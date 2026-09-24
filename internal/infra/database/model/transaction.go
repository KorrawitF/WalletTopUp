package model

import (
	"WalletTopUp/internal/domain/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const tenMinute = 10 * time.Minute

type Transaction struct {
	ID            string `gorm:"type:uuid;primaryKey"`
	UserID        uint
	Amount        float64
	PaymentMethod string
	Status        string
	ExpiresAt     time.Time
	CompletedAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (tx *Transaction) ToEntity() *entity.Transaction {
	return &entity.Transaction{
		ID:            tx.ID,
		UserID:        tx.UserID,
		Amount:        tx.Amount,
		PaymentMethod: entity.PaymentMethod(tx.PaymentMethod),
		Status:        entity.TransactionStatus(tx.Status),
		ExpiresAt:     tx.ExpiresAt,
	}
}

func (tx *Transaction) BeforeCreate(gTx *gorm.DB) error {
	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}

	tx.ID = uuid.String()
	tx.ExpiresAt = time.Now().Add(tenMinute)
	return nil
}
