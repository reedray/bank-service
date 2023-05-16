package grpc_transport

import (
	"context"
	"github.com/reedray/bank-service/api/pb/converter/gen_converter"
	"github.com/reedray/bank-service/internal/converter"
)

type Server struct {
	gen_converter.ConvertServiceServer
	useCase converter.ConvertUseCase
}

func (g *Server) Convert(ctx context.Context, money *gen_converter.Money) (*gen_converter.Money, error) {
	rate, err := g.useCase.Convert(ctx, money)
	if err != nil {
		return nil, err
	}
	return rate, nil
}

func (g *Server) mustEmbedUnimplementedConvertServiceServer() {
}

func NewServer(c converter.ConvertUseCase) *Server {
	return &Server{useCase: c}
}
