package core

import "errors"

var (
	ErrGeneric        = errors.New("something went wrong, check that the data in the geld.yaml is correct")
	ErrInvalidMode    = errors.New("invalid mode, the available modes are: mainnet, testnet, devnet")
	ErrPriceApiFailed = errors.New("something went wrong on the api")
)

type Wallet struct {
	Address      string
	Mnemonic     string
	PemPath      string
	JsonPath     string
	JsonPassword string
}

type MultiversxNetService interface {
	CreateWallet(folderPath string, password string) (Wallet, error)
}
