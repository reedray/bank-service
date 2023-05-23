package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/reedray/bank-service/internal/transact"
	"github.com/reedray/bank-service/internal/transact/entity"
)

type WithdrawUseCaseImpl struct {
	customerRepo transact.CustomerRepository
}

func NewWithdraw(c transact.CustomerRepository) *WithdrawUseCaseImpl {
	return &WithdrawUseCaseImpl{
		customerRepo: c,
	}
}

func (w *WithdrawUseCaseImpl) Execute(ctx context.Context, accountID uuid.UUID, amount float64, currCode string) error {
	customer, err := w.customerRepo.Find(ctx, accountID)
	if customer == nil || (customer.CreatedAt.IsZero() || err != nil) {
		return err
	}

	var total entity.CurrencyBalance
	json.Unmarshal(customer.BalanceRaw, &total)
	switch currCode {
	case "BYN":
		if total.BYN < amount {
			return fmt.Errorf("not enough funds to withdraw")
		}
		total.BYN -= amount
	case "USD":
		if total.USD < amount {
			return fmt.Errorf("not enough funds to withdraw")
		}
		total.USD -= amount
	case "EUR":
		if total.EUR < amount {
			return fmt.Errorf("not enough funds to withdraw")
		}
		total.EUR -= amount
	}
	bytes, err := json.Marshal(total)
	if err != nil {
		return err
	}
	customer.BalanceRaw = bytes
	err = w.customerRepo.Update(ctx, customer)
	if err != nil {
		return err
	}
	return nil
}
