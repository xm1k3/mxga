package controller

import (
	"log"

	"github.com/shopspring/decimal"
	"github.com/xm1k3/mxga/pkg/core"
	"github.com/xm1k3/mxga/pkg/core/services"
)

type ApiController struct {
	Service core.MultiversxApiService
}

type TokenFileData struct {
	Name       string `yaml:"name" json:"name"`
	Identifier string `yaml:"identifier" json:"identifier"`
}

func GetApiController(mode string) ApiController {
	tokenApiService := services.MultiversxApiService{
		ApiUrl: GetApiUrl(mode),
	}
	controller := ApiController{
		Service: tokenApiService,
	}
	return controller
}

func GetApiUrl(mode string) string {
	if mode == "" {
		return ""
	}
	switch mode {
	case "testnet":
		return "https://testnet-api.multiversx.com"
	case "devnet":
		return "https://devnet-api.multiversx.com"
	case "mainnet":
		return "https://api.multiversx.com"
	default:
		log.Fatal(core.ErrInvalidMode)
		return ""
	}
}

func GetTokenPrice(tokenIdentifier string, mode string) core.TokenApi {
	apiController := GetApiController(mode)
	token, err := apiController.Service.GetTokensPrice(tokenIdentifier)
	if err != nil {
		log.Fatal(err)
	}

	return token
}

func GetAccountTokenPrice(address string, tokenIdentifier string, mode string) decimal.Decimal {
	apiController := GetApiController(mode)
	price, err := apiController.Service.GetAccountToken(address, tokenIdentifier)
	if err != nil {
		log.Fatal(err)
	}

	return price
}
