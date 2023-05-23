package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/reedray/bank-service/internal/transact"
	"github.com/reedray/bank-service/internal/transact/entity"
)

type DepositUseCaseImpl struct {
	customerRepo transact.CustomerRepository
}

func NewDeposit(c transact.CustomerRepository) *DepositUseCaseImpl {
	return &DepositUseCaseImpl{
		customerRepo: c,
	}
}

func (d *DepositUseCaseImpl) Execute(ctx context.Context, accountID uuid.UUID, amount float64, currCode string) error {
	customer, err := d.customerRepo.Find(ctx, accountID)
	if customer == nil || (customer.CreatedAt.IsZero() || err != nil) {
		return err
	}

	var total entity.CurrencyBalance
	err = json.Unmarshal(customer.BalanceRaw, &total)
	if err != nil {
		fmt.Println(err)
		return err
	}
	switch currCode {
	case "BYN":
		total.BYN += amount
	case "USD":
		total.BYN += amount
	case "EUR":
		total.BYN += amount
	}
	bytes, err := json.Marshal(total)
	if err != nil {
		return err
	}
	customer.BalanceRaw = bytes
	err = d.customerRepo.Update(ctx, customer)
	if err != nil {
		return err
	}
	return nil

}
