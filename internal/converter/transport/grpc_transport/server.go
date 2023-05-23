package grpc_transport

import (
	"context"
	"github.com/reedray/bank-service/api/pb/converter/gen_convert"
	"github.com/reedray/bank-service/internal/converter"
)

type ConvertControllerImpl struct {
	gen_convert.UnsafeConvertServiceServer
	controller converter.ConvertController
}

func (g *ConvertControllerImpl) Convert(ctx context.Context, money *gen_convert.Money) (*gen_convert.Money, error) {
	rate, err := g.controller.Convert(ctx, money)
	if err != nil {
		return nil, err
	}
	return rate, nil
}

func NewConvertController(c converter.ConvertController) *ConvertControllerImpl {
	return &ConvertControllerImpl{controller: c}
}
