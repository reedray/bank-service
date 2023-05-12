package grpcHandler

import (
	"context"
	"fmt"
	"github.com/reedray/bank-service/app_services/converter/internal/usecase"
	"github.com/reedray/bank-service/app_services/converter/pkg/api/api_pb"
	"github.com/reedray/bank-service/pkg/logger"
	"strconv"
)

type GrpcHandler struct {
	api_pb.ConvertServiceServer
	useCase usecase.ConvertUseCase
	logger  logger.Logger
}

func (g *GrpcHandler) Convert(ctx context.Context, money *api_pb.Money) (*api_pb.Money, error) {
	dto := usecase.ExchangeDTO{CurrencyCode: money.CurrencyCode}
	rate, err := g.useCase.Convert(ctx, dto)
	if err != nil {
		g.logger.Error(err.Error())
		return nil, err
	}
	money.CurrencyCode = "BYN"
	amount, _ := strconv.ParseFloat(money.Amount, 64)
	amount *= rate
	money.Amount = fmt.Sprintf("%f", amount)
	return money, nil
}

func (g *GrpcHandler) mustEmbedUnimplementedConvertServiceServer() {
	//TODO implement me
	panic("implement me")
}

func New(c usecase.ConvertUseCase) *GrpcHandler {
	return &GrpcHandler{useCase: c}
}
