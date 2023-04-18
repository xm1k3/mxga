package core

import (
	"errors"

	"github.com/shopspring/decimal"
)

var (
	ErrGeneric        = errors.New("something went wrong, check that the data in the mxga.yaml is correct")
	ErrInvalidMode    = errors.New("invalid mode, the available modes are: mainnet, testnet, devnet")
	ErrPriceApiFailed = errors.New("something went wrong on the api")
)

type WalletConfig struct {
	Path  string
	Ext   string
	Main  string
	Other []string
}

type Wallet struct {
	Address      string
	Mnemonic     string
	PemPath      string
	JsonPath     string
	JsonPassword string
}

type MultiversxNetService interface {
	CreateWallet(folderPath string, password string) (Wallet, error)
	SendTransactions(pemPath string, to []string, amount decimal.Decimal, dataStr string) ([]string, error)

	GetAccount(address string) (string, error)
	GetTrxStatus(hash string) (string, error)
}
