package storage

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/reedray/bank-service/internal/transact/entity"
)

type TransactionRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewTransacioRepository(connPool *pgxpool.Pool) *TransactionRepositoryImpl {
	return &TransactionRepositoryImpl{db: connPool}
}

func (t *TransactionRepositoryImpl) Find(ctx context.Context, uuid uuid.UUID) (*entity.Transaction, error) {
	return nil, nil
}

func (t *TransactionRepositoryImpl) Save(ctx context.Context, tr *entity.Transaction) error {
	format := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	transaction := format.Insert("Transaction").Columns(
		"transactionID",
		"fromAccountID",
		"toAccountID",
		"amount",
		"currencyCode",
		"transactionTypeID",
	).Values(
		tr.TransactionID,
		tr.FromAccountID,
		tr.ToAccountID,
		tr.Amount,
		tr.CurrencyCode,
		tr.Type,
	)
	sql, args, err := transaction.ToSql()
	if err != nil {
		return err
	}

	_, err = t.db.Exec(ctx, sql, args...)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
