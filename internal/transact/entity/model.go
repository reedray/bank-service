package entity

import (
	"github.com/google/uuid"
	"time"
)

type Customer struct {
	ID         uuid.UUID `db:"id"`
	Username   string
	Password   string
	Role       string
	BalanceRaw []byte `json:"balance"`
	Status     string
	CreatedAt  time.Time
	DeletedAat time.Time
}

type Transaction struct {
	TransactionID uuid.UUID
	FromAccountID uuid.UUID
	ToAccountID   uuid.UUID
	Amount        float64
	CurrencyCode  string
	Type          int
}

type Loan struct {
	ID           uuid.UUID
	Amount       float64
	CurrencyCode string
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

type CurrencyBalance struct {
	BYN float64 `json:"BYN"`
	USD float64 `json:"USD"`
	EUR float64 `json:"EUR"`
}
