package usecase

import (
	"context"
	"github.com/reedray/bank-service/app_services/converter/internal/entity"
)

type ExchangeDTO struct {
	CurrencyCode string
}

type ConvertUseCase interface {
	Convert(ctx context.Context, data ExchangeDTO) (float64, error)
}

type ConvertRepository interface {
	GetExchangeRates(context.Context, string) (float64, error)
	SetExchangeRates(context.Context, string, float64) error
}

type ConvertWebAPI interface {
	Convert(int) (entity.ExchangeRate, error)
}
