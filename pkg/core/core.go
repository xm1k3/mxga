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
	Path    string
	Ext     string
	Main    string
	Address string
	Other   []string
}

type Wallet struct {
	Address      string
	Mnemonic     string
	PemPath      string
	JsonPath     string
	JsonPassword string
}

type TokenApi struct {
	Identifier        string          `json:"identifier"`
	Name              string          `json:"name"`
	Ticker            string          `json:"ticker"`
	Owner             string          `json:"owner"`
	Minted            string          `json:"minted"`
	Burnt             string          `json:"burnt"`
	InitialMinted     string          `json:"initialMinted"`
	Decimals          int             `json:"decimals"`
	IsPaused          bool            `json:"isPaused"`
	Transactions      int             `json:"transactions"`
	Accounts          int             `json:"accounts"`
	CanUpgrade        bool            `json:"canUpgrade"`
	CanMint           bool            `json:"canMint"`
	CanBurn           bool            `json:"canBurn"`
	CanChangeOwner    bool            `json:"canChangeOwner"`
	CanPause          bool            `json:"canPause"`
	CanFreeze         bool            `json:"canFreeze"`
	CanWipe           bool            `json:"canWipe"`
	Price             decimal.Decimal `json:"price"`
	MarketCap         decimal.Decimal `json:"marketCap"`
	Supply            string          `json:"supply"`
	CirculatingSupply string          `json:"circulatingSupply"`
}

type SwapResult struct {
	Status string `json:"status"`
	Hash   string `json:"hash"`
}

type MultiversxNetService interface {
	CreateWallet(folderPath string, password string) (Wallet, error)
	SendTransactions(pemPath string, to []string, amount decimal.Decimal, dataStr string) ([]string, error)
	CreateSwapTokensFixedInput(pemPath string, contract string, fromToken string, amount decimal.Decimal, toToken string, slippage decimal.Decimal, fromDecimals int, toDecimals int) (SwapResult, error)

	GetAccount(address string) (string, error)
	GetTrxStatus(hash string) (string, error)
}

type MultiversxApiService interface {
	GetTokensPrice(tokens string) (TokenApi, error)
	GetAccountToken(address string, token string) (decimal.Decimal, error)
}
