package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/shopspring/decimal"
	"github.com/xm1k3/mxga/pkg/core"
	"github.com/xm1k3/mxga/pkg/utils"
)

type MultiversxApiService struct {
	ApiUrl string
}

func (m MultiversxApiService) GetTokensPrice(token string) (core.TokenApi, error) {
	req, err := http.NewRequest("GET", m.ApiUrl+"/tokens/"+token, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return core.TokenApi{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return core.TokenApi{}, core.ErrPriceApiFailed
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var price core.TokenApi
	err = json.Unmarshal([]byte(body), &price)
	if err != nil {
		return core.TokenApi{}, err
	}
	return price, err
}

func (m MultiversxApiService) GetAccountToken(address string, token string) (decimal.Decimal, error) {

	type AccountStruct struct {
		Identifier        string          `json:"identifier"`
		Name              string          `json:"name"`
		Ticker            string          `json:"ticker"`
		CirculatingSupply string          `json:"circulatingSupply"`
		Balance           string          `json:"balance"`
		Decimals          int             `json:"decimals"`
		ValueUsd          decimal.Decimal `json:"valueUsd"`
	}

	req, err := http.NewRequest("GET", m.ApiUrl+"/accounts/"+address+"/tokens/"+token, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return decimal.Decimal{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return decimal.Decimal{}, core.ErrPriceApiFailed
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var accountToken AccountStruct
	err = json.Unmarshal([]byte(body), &accountToken)
	if err != nil {
		return decimal.Decimal{}, err
	}

	decimalValueStr := utils.ConvertStringToDecimal(accountToken.Balance, accountToken.Decimals)

	decimalValue, err := decimal.NewFromString(decimalValueStr)
	if err != nil {
		return decimal.Decimal{}, err
	}
	decimalValue = decimalValue.Sub(decimal.NewFromFloat(0.02))

	return decimalValue, err
}
