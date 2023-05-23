package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/reedray/bank-service/internal/transact"
	"github.com/reedray/bank-service/internal/transact/entity"
)

type TransferUseCaseImpl struct {
	customerRepo transact.CustomerRepository
}

func NewTransfer(c transact.CustomerRepository) *TransferUseCaseImpl {
	return &TransferUseCaseImpl{
		customerRepo: c,
	}
}

func (t *TransferUseCaseImpl) Execute(ctx context.Context, fromAccountID, toAccountID uuid.UUID, amount float64, currCode string) error {
	from, err := t.customerRepo.Find(ctx, fromAccountID)
	if from == nil || (from.CreatedAt.IsZero() || err != nil) {
		return err
	}
	to, err := t.customerRepo.Find(ctx, toAccountID)
	if to == nil || (to.CreatedAt.IsZero() || err != nil) {
		return err
	}

	var totalFrom entity.CurrencyBalance
	var totalTo entity.CurrencyBalance
	var tf float64
	json.Unmarshal(from.BalanceRaw, &totalFrom)
	json.Unmarshal(to.BalanceRaw, &totalTo)
	switch currCode {
	case "BYN":
		tf = totalFrom.BYN
	case "USD":
		tf = totalFrom.USD
	case "EUR":
		tf = totalFrom.EUR
	}
	if amount > tf {
		return fmt.Errorf("not enough funds to transfer")
	}
	switch currCode {
	case "BYN":
		totalFrom.BYN -= amount
		totalTo.BYN += amount
	case "USD":
		totalFrom.USD -= amount
		totalTo.USD += amount
	case "EUR":
		totalFrom.EUR -= amount
		totalTo.EUR += amount
	}

	bytes, err := json.Marshal(totalFrom)
	if err != nil {
		return err
	}
	from.BalanceRaw = bytes
	bytesTo, err := json.Marshal(totalTo)
	if err != nil {
		return err
	}
	to.BalanceRaw = bytesTo
	err = t.customerRepo.Update(ctx, from)
	if err != nil {
		return err
	}
	err = t.customerRepo.Update(ctx, to)
	if err != nil {
		return err
	}
	return nil
}
