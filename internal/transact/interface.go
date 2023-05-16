package transact

import (
	"github.com/google/uuid"
	"github.com/reedray/bank-service/internal/transact/entity"
)

type TransactService interface {
	GetBalance() (float64, error)
}

type TransactUseCase interface {
	GetBalance()
	Deposit()
	Withdraw()
	Transfer()
}

type AccountRepository interface {
	GetBalance(userId uuid.UUID) (entity.Balance, error)
	Deposit(userID uuid.UUID) (float64, error)
	Withdraw(userID uuid.UUID) (float64, error)
	Transfer(amount float64, fromID, toID uuid.UUID) error
}
