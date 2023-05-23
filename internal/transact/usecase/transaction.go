package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/reedray/bank-service/internal/transact"
	"github.com/reedray/bank-service/internal/transact/entity"
)

type TransactionUseCaseImpl struct {
	transactionRepo transact.TransactionRepository
	db              *pgxpool.Pool
}

func NewTransaction(tr transact.TransactionRepository, db *pgxpool.Pool) *TransactionUseCaseImpl {
	return &TransactionUseCaseImpl{transactionRepo: tr, db: db}
}

func (t *TransactionUseCaseImpl) Begin(
	ctx context.Context,
	fromID uuid.UUID,
	toId uuid.UUID,
	amount float64,
	currCode string,
	trType int,
) (pgx.Tx, *entity.Transaction, error) {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	transaction := entity.Transaction{
		TransactionID: uuid.New(),
		FromAccountID: fromID,
		ToAccountID:   toId,
		Amount:        amount,
		CurrencyCode:  currCode,
		Type:          trType,
	}
	return tx, &transaction, nil
}

func (t *TransactionUseCaseImpl) Rollback(ctx context.Context, tx pgx.Tx) error {
	err := tx.Rollback(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (t *TransactionUseCaseImpl) Commit(ctx context.Context, tx pgx.Tx, etx *entity.Transaction) error {

	err := t.transactionRepo.Save(ctx, etx)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
