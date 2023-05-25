package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/reedray/bank-service/internal/transact"
	"log"
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
	log.Println("Looking for customer in bd")
	customer, err := b.customerRepo.Find(ctx, accountID)
	if customer == nil || (customer.CreatedAt.IsZero() || err != nil) {
		log.Println("Failed to find customer", err)
		return nil, err
	}
	return customer.BalanceRaw, nil

}
