package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/reedray/bank-service/internal/transact"
)

type BalanceUseCaseImpl struct {
	customerRepo transact.CustomerRepository
}

func NewBalance(c transact.CustomerRepository) *BalanceUseCaseImpl {
	return &BalanceUseCaseImpl{
		customerRepo: c,
	}
}

func (b *BalanceUseCaseImpl) Execute(ctx context.Context, accountID uuid.UUID) ([]byte, error) {
	customer, err := b.customerRepo.Find(ctx, accountID)
	if customer == nil || (customer.CreatedAt.IsZero() || err != nil) {
		return nil, err
	}
	return customer.BalanceRaw, nil

}
