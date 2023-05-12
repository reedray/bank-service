package usecase

import (
	"context"
	"fmt"
)

// for API of national bank
var hm = map[string]int{"USD": 431, "EUR": 451}

type ConvertUseCaseImpl struct {
	repo     ConvertRepository
	webApi   ConvertWebAPI
	Codes_ID map[string]int
}

func New(cr ConvertRepository, wa ConvertWebAPI) *ConvertUseCaseImpl {
	return &ConvertUseCaseImpl{
		repo:     cr,
		webApi:   wa,
		Codes_ID: hm,
	}
}

func (c *ConvertUseCaseImpl) Convert(ctx context.Context, rate ExchangeDTO) (float64, error) {
	repoRate, err := c.repo.GetExchangeRates(ctx, rate.CurrencyCode)
	if err != nil {
		webApiResponse, err := c.webApi.Convert(c.Codes_ID[rate.CurrencyCode])
		if err != nil {
			return 0, fmt.Errorf("can`t pull repoRate %w", err)
		}
		//setting repoRate
		err = c.repo.SetExchangeRates(ctx, webApiResponse.CurrencyCode, webApiResponse.Rate)
		if err != nil {
			//todo: log message
		}
		return webApiResponse.Rate, nil
	}
	return repoRate, nil
}
