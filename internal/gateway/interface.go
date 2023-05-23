package gateway

import (
	"context"
	"github.com/reedray/bank-service/api/pb/transact/gen_transact"
)

type GatewayService interface {
	Login(ctx context.Context, request *gen_transact.LoginRequest) *gen_transact.LoginResponse
	Register(ctx context.Context, request *gen_transact.RegisterRequest) *gen_transact.RegisterResponse
	Deposit(ctx context.Context, request *gen_transact.DepositRequest) *gen_transact.Error
	Withdraw(context.Context, *gen_transact.WithdrawRequest) *gen_transact.Error
	Transfer(context.Context, *gen_transact.TransferRequest) *gen_transact.Error
	Balance(context.Context, *gen_transact.BalanceRequest) *gen_transact.BalanceResponse
}
