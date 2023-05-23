package transact

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/reedray/bank-service/api/pb/transact/gen_transact"
	"github.com/reedray/bank-service/internal/transact/entity"
)

type TransferHandler interface {
	Login(ctx context.Context, request *gen_transact.LoginRequest) *gen_transact.LoginResponse
	Register(ctx context.Context, request *gen_transact.RegisterRequest) *gen_transact.RegisterResponse
	Deposit(ctx context.Context, request *gen_transact.DepositRequest) *gen_transact.Error
	Withdraw(context.Context, *gen_transact.WithdrawRequest) *gen_transact.Error
	Transfer(context.Context, *gen_transact.TransferRequest) *gen_transact.Error
	Balance(context.Context, *gen_transact.BalanceRequest) *gen_transact.BalanceResponse
}

type DepositUseCase interface {
	Execute(ctx context.Context, accountID uuid.UUID, amount float64, currCode string) error
}

type WithdrawUseCase interface {
	Execute(ctx context.Context, accountID uuid.UUID, amount float64, currCode string) error
}

type TransferUseCase interface {
	Execute(ctx context.Context, fromAccountID, toAccountID uuid.UUID, amount float64, currCode string) error
}

type BalanceUseCase interface {
	Execute(ctx context.Context, accountID uuid.UUID) ([]byte, error)
}

type CustomerRepository interface {
	Find(ctx context.Context, accountID uuid.UUID) (*entity.Customer, error)
	Save(ctx context.Context, account *entity.Customer) error
	Update(ctx context.Context, account *entity.Customer) error
}

type TransactionUseCase interface {
	Begin(ctx context.Context, fromID uuid.UUID, toID uuid.UUID, amount float64, currCode string, trType int) (pgx.Tx, *entity.Transaction, error)
	Rollback(ctx context.Context, tx pgx.Tx) error
	Commit(ctx context.Context, tx pgx.Tx, etx *entity.Transaction) error
}

type TransactionRepository interface {
	Find(context.Context, uuid.UUID) (*entity.Transaction, error)
	Save(context.Context, *entity.Transaction) error
}

type AuthUseCase interface {
	Login(ctx context.Context, username, password, secret string) (string, error)
	Register(ctx context.Context, username, password, secret string) (string, error)
	ValidateToken(tokenString, secret string) (bool, error)
}
type AuthRepository interface {
	FindByCredentials(ctx context.Context, username, password string) (*entity.Customer, error)
}
