package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/reedray/bank-service/api/pb/transact/gen_transact"
	cfg "github.com/reedray/bank-service/config/transact"
	"github.com/reedray/bank-service/internal/transact"
	"github.com/reedray/bank-service/internal/transact/entity"
	"strconv"
)

const (
	withdrawType = 1
	depositType  = 2
	transferType = 3
	balanceType  = 4
)

type TransferHandlerImpl struct {
	cfg *cfg.Config
	transact.AuthUseCase
	transact.DepositUseCase
	transact.WithdrawUseCase
	transact.TransferUseCase
	transact.BalanceUseCase
	transact.TransactionUseCase
}

func NewTransferHandler(
	cfg *cfg.Config,
	a transact.AuthUseCase,
	d transact.DepositUseCase,
	w transact.WithdrawUseCase,
	tr transact.TransferUseCase,
	b transact.BalanceUseCase,
	tx transact.TransactionUseCase,
) *TransferHandlerImpl {
	return &TransferHandlerImpl{
		cfg:                cfg,
		AuthUseCase:        a,
		DepositUseCase:     d,
		WithdrawUseCase:    w,
		TransferUseCase:    tr,
		BalanceUseCase:     b,
		TransactionUseCase: tx,
	}
}

func (t *TransferHandlerImpl) toUUID(id string) (uuid.UUID, error) {
	parse, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, err
	}
	return parse, err
}
func (t *TransferHandlerImpl) toFloat(s string) float64 {
	float, _ := strconv.ParseFloat(s, 64)
	return float
}
func (t *TransferHandlerImpl) extractClaims(tokenString string) (jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		return []byte(t.cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	// Check if the token is valid and contains claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token or missing claims")
}
func (t *TransferHandlerImpl) getFromClaims(claims jwt.MapClaims, key string) (string, error) {
	if username, ok := claims[key].(string); ok {
		return username, nil
	}
	return "", fmt.Errorf("%s claim not found or not a string", key)
}

func (t *TransferHandlerImpl) Login(ctx context.Context, request *gen_transact.LoginRequest) *gen_transact.LoginResponse {

	token, err := t.AuthUseCase.Login(ctx, request.Username, request.Password, t.cfg.Secret)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &gen_transact.LoginResponse{Token: token}
}

func (t *TransferHandlerImpl) Register(ctx context.Context, request *gen_transact.RegisterRequest) *gen_transact.RegisterResponse {
	token, err := t.AuthUseCase.Register(ctx, request.Username, request.Password, t.cfg.Secret)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &gen_transact.RegisterResponse{Token: token}
}

func (t *TransferHandlerImpl) Deposit(ctx context.Context, request *gen_transact.DepositRequest) *gen_transact.Error {
	resp := &gen_transact.Error{}
	valid, err := t.AuthUseCase.ValidateToken(request.Token, t.cfg.Secret)
	if err != nil || !valid {
		resp.Err = err.Error()
		return resp
	}

	claims, err := t.extractClaims(request.Token)
	if err != nil {
		resp.Err = err.Error()
		return resp
	}
	id, err := t.getFromClaims(claims, "id")
	if err != nil {
		resp.Err = err.Error()
		return resp
	}

	toUUID, err := t.toUUID(id)
	if err != nil {
		resp.Err = err.Error()
		return resp
	}

	amount := t.toFloat(request.Total.Amount)

	tx, transaction, err := t.TransactionUseCase.Begin(ctx, toUUID, toUUID, amount, request.Total.CurrencyCode, depositType)
	if err != nil {
		resp.Err = err.Error()
		return resp
	}

	err = t.DepositUseCase.Execute(ctx, toUUID, amount, request.Total.CurrencyCode)
	if err != nil {
		err1 := t.TransactionUseCase.Rollback(ctx, tx)
		if err1 != nil {
			resp.Err = err.Error()
			return resp
		}
		resp.Err = err.Error()
		return resp
	}
	err = t.TransactionUseCase.Commit(ctx, tx, transaction)
	if err != nil {

		resp.Err = err.Error()
		return resp
	}
	return nil
}

func (t *TransferHandlerImpl) Withdraw(ctx context.Context, request *gen_transact.WithdrawRequest) *gen_transact.Error {
	resp := &gen_transact.Error{}
	valid, err := t.AuthUseCase.ValidateToken(request.Token, t.cfg.Secret)
	if err != nil || !valid {
		resp.Err = err.Error()
		return resp
	}

	claims, err := t.extractClaims(request.Token)
	if err != nil {
		resp.Err = err.Error()
		return resp
	}
	id, err := t.getFromClaims(claims, "id")
	if err != nil {
		return nil
	}

	toUUID, err := t.toUUID(id)
	if err != nil {
		resp.Err = err.Error()
		return resp
	}

	amount := t.toFloat(request.Total.Amount)

	tx, transaction, err := t.TransactionUseCase.Begin(ctx, toUUID, toUUID, amount, request.Total.CurrencyCode, withdrawType)
	if err != nil {
		resp.Err = err.Error()
		return resp
	}

	err = t.WithdrawUseCase.Execute(ctx, toUUID, amount, request.Total.CurrencyCode)
	if err != nil {
		err1 := t.TransactionUseCase.Rollback(ctx, tx)
		if err1 != nil {
			resp.Err = err1.Error()
			return resp
		}
		resp.Err = err.Error()
		return resp
	}
	err = t.TransactionUseCase.Commit(ctx, tx, transaction)
	if err != nil {

		resp.Err = err.Error()
		return resp
	}
	return nil
}

func (t *TransferHandlerImpl) Transfer(ctx context.Context, request *gen_transact.TransferRequest) *gen_transact.Error {
	resp := &gen_transact.Error{}
	valid, err := t.AuthUseCase.ValidateToken(request.Token, t.cfg.Secret)
	if err != nil || !valid {
		resp.Err = err.Error()
		return resp
	}

	claims, err := t.extractClaims(request.Token)
	if err != nil {
		resp.Err = err.Error()
		return resp
	}
	idStr, err := t.getFromClaims(claims, "id")
	if err != nil {
		return nil
	}

	id, err := t.toUUID(idStr)
	if err != nil {
		resp.Err = err.Error()
		return resp
	}

	amount := t.toFloat(request.Total.Amount)

	toID, err := t.toUUID(request.IdTo)
	if err != nil {
		resp.Err = err.Error()
		return resp
	}
	tx, transaction, err := t.TransactionUseCase.Begin(ctx, id, toID, amount, request.Total.CurrencyCode, transferType)
	if err != nil {
		resp.Err = err.Error()
		return resp
	}

	err = t.TransferUseCase.Execute(ctx, id, toID, amount, request.Total.CurrencyCode)
	if err != nil {
		err1 := t.TransactionUseCase.Rollback(ctx, tx)
		if err1 != nil {
			resp.Err = err.Error()
			return resp
		}
		resp.Err = err.Error()
		return resp
	}
	err = t.TransactionUseCase.Commit(ctx, tx, transaction)
	if err != nil {
		resp.Err = err.Error()
		return resp
	}
	return nil
}

func (t *TransferHandlerImpl) Balance(ctx context.Context, request *gen_transact.BalanceRequest) *gen_transact.BalanceResponse {
	resp := &gen_transact.BalanceResponse{}
	valid, err := t.AuthUseCase.ValidateToken(request.Token, t.cfg.Secret)
	if err != nil || !valid {
		resp.ErrorMsg.Err = err.Error()
		return resp
	}

	claims, err := t.extractClaims(request.Token)
	if err != nil {
		resp.ErrorMsg.Err = err.Error()
		return resp
	}
	id, err := t.getFromClaims(claims, "id")
	if err != nil {
		resp.ErrorMsg.Err = err.Error()
		return resp
	}

	toUUID, err := t.toUUID(id)
	if err != nil {
		resp.ErrorMsg.Err = err.Error()
		return resp
	}

	tx, transaction, err := t.TransactionUseCase.Begin(ctx, toUUID, toUUID, 0, "", balanceType)
	if err != nil {
		resp.ErrorMsg.Err = err.Error()
		return resp
	}

	bytes, err := t.BalanceUseCase.Execute(ctx, toUUID)
	if err != nil {
		err = t.TransactionUseCase.Rollback(ctx, tx)
		if err != nil {
			resp.ErrorMsg.Err = err.Error()
			return resp
		}
		resp.ErrorMsg.Err = err.Error()
		return resp
	}
	err = t.TransactionUseCase.Commit(ctx, tx, transaction)
	if err != nil {

		resp.ErrorMsg.Err = err.Error()
		return resp
	}
	var money entity.CurrencyBalance
	json.Unmarshal(bytes, &money)
	byn := strconv.FormatFloat(money.BYN, 'f', -1, 64)
	usd := strconv.FormatFloat(money.USD, 'f', -1, 64)
	eur := strconv.FormatFloat(money.EUR, 'f', -1, 64)
	resp.Total.Byn = byn
	resp.Total.Usd = usd
	resp.Total.Eur = eur
	return resp
}
