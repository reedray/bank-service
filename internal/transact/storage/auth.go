package storage

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/reedray/bank-service/internal/transact/entity"
	"log"
)

type AuthRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewAuthRepositoty(db *pgxpool.Pool) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{db: db}
}

func (a *AuthRepositoryImpl) FindByCredentials(ctx context.Context, username, password string) (*entity.Customer, error) {
	format := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	users := format.Select("*").From("Customer").
		Where(sq.Eq{"username": username}).
		Where(sq.Eq{"password": password})
	sql, args, err := users.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := a.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res entity.Customer
	for rows.Next() {
		err := rows.Scan(&res.ID, &res.Username, &res.Password, &res.Role, &res.CreatedAt, &res.DeletedAat, &res.Status, &res.BalanceRaw)
		if err != nil {
			return nil, err
		}
	}
	if rows.Err() != nil {
		log.Printf("failed to iterate over rows: %v", rows.Err())
	}
	return &res, nil
}
