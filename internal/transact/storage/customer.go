package storage

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/reedray/bank-service/internal/transact/entity"
	"log"
)

type CustomerRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewCustomerRepository(pool *pgxpool.Pool) *CustomerRepositoryImpl {
	return &CustomerRepositoryImpl{db: pool}
}

func (c *CustomerRepositoryImpl) Find(ctx context.Context, accountID uuid.UUID) (*entity.Customer, error) {
	log.Println("Finding user with id", accountID)
	format := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	user := format.Select("*").From("Customer").Where(sq.Eq{"id": accountID})
	sql, args, err := user.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := c.db.Query(ctx, sql, args...)
	if err != nil {
		log.Println("failed to fetch data with error:", err)
		return nil, err
	}
	defer rows.Close()
	res := entity.Customer{}
	for rows.Next() {
		err1 := rows.Scan(&res.ID, &res.Username, &res.Password, &res.Role, &res.CreatedAt, &res.DeletedAat, &res.Status, &res.BalanceRaw)
		if err1 != nil {
			log.Println("failed to fetch data with error:", err1)
			return nil, err1
		}
	}
	if rows.Err() != nil {
		log.Printf("failed to iterate over rows: %v", rows.Err())
	}
	return &res, nil
}

func (c *CustomerRepositoryImpl) Save(ctx context.Context, account *entity.Customer) error {
	format := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := format.Insert("Customer").Columns(
		"id",
		"username",
		"password",
		"role",
		"createdAt",
		"deletedAt",
		"status",
		"balance",
	).Values(
		account.ID,
		account.Username,
		account.Password,
		account.Role,
		account.CreatedAt,
		account.DeletedAat,
		account.Status,
		account.BalanceRaw,
	).ToSql()
	if err != nil {
		return err
	}
	_, err = c.db.Exec(ctx, sql, args...)
	if err != nil {
		fmt.Println("internal db error")
		return err
	}
	return nil
}

func (c *CustomerRepositoryImpl) Update(ctx context.Context, account *entity.Customer) error {
	log.Println("Updating user", account.ID)
	format := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	data := map[string]interface{}{
		"username":  account.Username,
		"password":  account.Password,
		"role":      account.Role,
		"createdAt": account.CreatedAt,
		"deletedAt": account.DeletedAat,
		"status":    account.Status,
		"balance":   account.BalanceRaw,
	}
	sql, args, err := format.Update("Customer").
		SetMap(data).
		Where(sq.Eq{"username": account.Username, "password": account.Password}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = c.db.Exec(ctx, sql, args...)
	if err != nil {
		log.Printf("Failed to update user %s with error %s", account.ID.String(), err.Error())
		return err
	}
	return nil
}
