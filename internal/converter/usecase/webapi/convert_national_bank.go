package webapi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/reedray/bank-service/internal/converter/entity"
	"net/http"
	"strconv"
)

const (
	url = "https://api.nbrb.by/exrates/rates/"
)

type responseWebAPI struct {
	CurID           int     `json:"Cur_ID"`
	Date            string  `json:"Date"`
	CurAbbreviation string  `json:"Cur_Abbreviation"`
	CurScale        int     `json:"Cur_Scale"`
	CurName         string  `json:"Cur_Name"`
	CurOfficialRate float64 `json:"Cur_OfficialRate"`
}

type ConvertWebAPI struct {
	client *http.Client
}

func New() *ConvertWebAPI {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Set this to true if you want to skip SSL certificate verification !!!only for testing
			},
		},
	}
	return &ConvertWebAPI{client: client}
}

func (cw *ConvertWebAPI) Convert(currencyId int) (entity.ExchangeRate, error) {
	resp, err := cw.client.Get(url + strconv.Itoa(currencyId))
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("web api failed to fetch data")
		return entity.ExchangeRate{}, fmt.Errorf("web api failed to fetch data %w", err)
	}
	defer resp.Body.Close()

	r := responseWebAPI{}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return entity.ExchangeRate{}, fmt.Errorf("can`t decode respone body %w", err)
	}

	return entity.ExchangeRate{
		CurrencyID:   r.CurID,
		CurrencyCode: r.CurAbbreviation,
		Rate:         r.CurOfficialRate,
	}, nil

}
