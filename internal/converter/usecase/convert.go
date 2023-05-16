package usecase

import (
	"context"
	"fmt"
	"github.com/reedray/bank-service/api/pb/converter/gen_converter"
	"github.com/reedray/bank-service/internal/converter"
	"strconv"
)

// for API of national bank
var hm = map[string]int{"USD": 431, "EUR": 451}

type ConvertUseCaseImpl struct {
	repo     converter.ConvertRepository
	webApi   converter.ConvertWebAPI
	Codes_ID map[string]int
}

func New(cr converter.ConvertRepository, wa converter.ConvertWebAPI) *ConvertUseCaseImpl {
	return &ConvertUseCaseImpl{
		repo:     cr,
		webApi:   wa,
		Codes_ID: hm,
	}
}

func (c *ConvertUseCaseImpl) Convert(ctx context.Context, data *gen_converter.Money) (*gen_converter.Money, error) {

	repoRate, err := c.repo.GetExchangeRates(ctx, data.CurrencyCode)
	if err != nil {
		webApiResponse, err := c.webApi.Convert(c.Codes_ID[data.CurrencyCode]) //check existence
		if err != nil {
			return nil, err
		}
		//setting repoRate
		err = c.repo.SetExchangeRates(ctx, webApiResponse.CurrencyCode, webApiResponse.Rate)
		if err != nil {
			//todo: log message
		}
		repoRate = webApiResponse.Rate
	}
	//possibly should be as a method of entity layer
	m := gen_converter.Money{
		Amount:       "",
		CurrencyCode: "BYN",
	}
	amount, _ := strconv.ParseFloat(data.Amount, 64)
	amount *= repoRate
	m.Amount = fmt.Sprintf("%f", amount)
	return &m, nil
}
