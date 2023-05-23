package grpc_transport

import (
	"context"
	"fmt"
	"github.com/reedray/bank-service/api/pb/transact/gen_transact"
	"github.com/reedray/bank-service/internal/transact"
)

type ServerHandler struct {
	transact.TransferHandler
	gen_transact.UnimplementedTransferServiceServer
}

func NewServerHandler(h transact.TransferHandler) *ServerHandler {
	return &ServerHandler{
		TransferHandler: h,
	}
}

func (s *ServerHandler) Transfer(ctx context.Context, request *gen_transact.TransferRequest) (*gen_transact.Error, error) {
	resp := s.TransferHandler.Transfer(ctx, request)
	if resp.Err != "" {
		return resp, fmt.Errorf(resp.Err)
	}
	return nil, nil
}

func (s *ServerHandler) Withdraw(ctx context.Context, request *gen_transact.WithdrawRequest) (*gen_transact.Error, error) {
	resp := s.TransferHandler.Withdraw(ctx, request)
	if resp.Err != "" {
		return resp, fmt.Errorf(resp.Err)
	}
	return nil, nil
}

func (s *ServerHandler) Deposit(ctx context.Context, request *gen_transact.DepositRequest) (*gen_transact.Error, error) {
	resp := s.TransferHandler.Deposit(ctx, request)
	if resp.Err != "" {
		return resp, fmt.Errorf(resp.Err)
	}
	return nil, nil
}

func (s *ServerHandler) Balance(ctx context.Context, request *gen_transact.BalanceRequest) (*gen_transact.BalanceResponse, error) {
	resp := s.TransferHandler.Balance(ctx, request)
	if resp.ErrorMsg.Err != "" {
		return resp, fmt.Errorf(resp.ErrorMsg.Err)
	}
	return resp, nil
}

func (s *ServerHandler) Login(ctx context.Context, request *gen_transact.LoginRequest) (*gen_transact.LoginResponse, error) {
	resp := s.TransferHandler.Login(ctx, request)
	return resp, nil
}

func (s *ServerHandler) Register(ctx context.Context, request *gen_transact.RegisterRequest) (*gen_transact.RegisterResponse, error) {
	resp := s.TransferHandler.Register(ctx, request)
	return resp, nil
}
