package converter

import (
	"context"
	"github.com/reedray/bank-service/api/pb/converter/gen_convert"
	"github.com/reedray/bank-service/internal/converter/entity"
)

type ConvertController interface {
	Convert(ctx context.Context, data *gen_convert.Money) (*gen_convert.Money, error)
}

type ConvertRepository interface {
	GetExchangeRates(context.Context, string) (float64, error)
	SetExchangeRates(context.Context, string, float64) error
}

type ConvertUseCase interface {
	Convert(int) (entity.ExchangeRate, error)
}
