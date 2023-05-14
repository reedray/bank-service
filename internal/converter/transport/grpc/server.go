package grpc

import (
	"context"
	"github.com/reedray/bank-service/api/pb/converter"
	converter2 "github.com/reedray/bank-service/internal/converter"
)

type Server struct {
	converter.ConvertServiceServer
	useCase converter2.ConvertUseCase
}

func (g *Server) Convert(ctx context.Context, money *converter.Money) (*converter.Money, error) {
	rate, err := g.useCase.Convert(ctx, money)
	if err != nil {
		return nil, err
	}
	return rate, nil
}

func (g *Server) mustEmbedUnimplementedConvertServiceServer() {
}

func New(c converter2.ConvertUseCase) *Server {
	return &Server{useCase: c}
}
