package entity

import (
	"time"
)

type TransactionStatus string

const (
	TransactionStatusVerified  TransactionStatus = "verified"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusExpired   TransactionStatus = "expired"
)

type PaymentMethod string

const (
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodDebitCard    PaymentMethod = "debit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
)

func (p PaymentMethod) IsValid() bool {
	switch p {
	case PaymentMethodCreditCard, PaymentMethodDebitCard, PaymentMethodBankTransfer:
		return true
	}
	return false
}

type Transaction struct {
	ID            string
	UserID        uint
	Amount        float64
	PaymentMethod PaymentMethod
	Status        TransactionStatus
	ExpiresAt     time.Time
	CompletedAt   *time.Time
}

func (tx *Transaction) IsExpired() bool {
	now := time.Now()
	return !tx.ExpiresAt.Before(now)
}
