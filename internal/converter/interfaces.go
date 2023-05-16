package converter

import (
	"context"
	"github.com/reedray/bank-service/api/pb/converter/gen_converter"
	"github.com/reedray/bank-service/internal/converter/entity"
)

type ConvertUseCase interface {
	Convert(ctx context.Context, data *gen_converter.Money) (*gen_converter.Money, error)
}

type ConvertRepository interface {
	GetExchangeRates(context.Context, string) (float64, error)
	SetExchangeRates(context.Context, string, float64) error
}

type ConvertWebAPI interface {
	Convert(int) (entity.ExchangeRate, error)
}
