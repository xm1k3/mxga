package services

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/multiversx/mx-chain-crypto-go/signing"
	"github.com/multiversx/mx-chain-crypto-go/signing/ed25519"
	"github.com/multiversx/mx-sdk-go/blockchain"
	"github.com/multiversx/mx-sdk-go/blockchain/cryptoProvider"
	"github.com/multiversx/mx-sdk-go/builders"
	mxcore "github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/interactors"
	"github.com/shopspring/decimal"
	"github.com/xm1k3/mxga/pkg/core"
	"github.com/xm1k3/mxga/pkg/utils"
)

type MultiversxService struct {
	ProxyUrl string
}

var (
	suite  = ed25519.NewEd25519()
	keyGen = signing.NewKeyGenerator(suite)
)

// Function for create a single wallet, returning a struct
func (m MultiversxService) CreateWallet(folderPath string, password string) (core.Wallet, error) {
	exPath := filepath.Dir(folderPath)
	exPath = path.Join(exPath, "wallets")

	w := interactors.NewWallet()
	mnemonic, err := w.GenerateMnemonic()
	if err != nil {
		return core.Wallet{}, err
	}

	walletPrivateKey := w.GetPrivateKeyFromMnemonic(mnemonic, 0, 0)
	wallet, err := w.GetAddressFromPrivateKey(walletPrivateKey)
	if err != nil {
		return core.Wallet{}, err
	}

	// Create folder if not exists
	err = os.MkdirAll(exPath, os.ModePerm)
	if err != nil {
		return core.Wallet{}, err
	}

	// Pem file path
	err = os.MkdirAll(path.Join(exPath, "pem"), os.ModePerm)
	if err != nil {
		return core.Wallet{}, err
	}

	// Json file path
	err = os.MkdirAll(path.Join(exPath, "json"), os.ModePerm)
	if err != nil {
		return core.Wallet{}, err
	}

	// Save Pem file
	pemPrivatePath := path.Join(exPath, "pem", wallet.AddressAsBech32String()+".pem")
	err = w.SavePrivateKeyToPemFile(walletPrivateKey, pemPrivatePath)
	if err != nil {
		return core.Wallet{}, err
	}

	// Save Json file
	jsonPrivatePath := path.Join(exPath, "json", wallet.AddressAsBech32String()+".json")
	err = w.SavePrivateKeyToJsonFile(walletPrivateKey, password, jsonPrivatePath)
	if err != nil {
		return core.Wallet{}, err
	}
	return core.Wallet{
		Address:      wallet.AddressAsBech32String(),
		Mnemonic:     string(mnemonic),
		PemPath:      pemPrivatePath,
		JsonPath:     jsonPrivatePath,
		JsonPassword: password,
	}, nil
}

func (m MultiversxService) SendTransactions(pemPath string, to []string, amount decimal.Decimal, dataStr string) ([]string, error) {
	var toAddresses []string
	for _, t := range to {
		addr, err := data.NewAddressFromBech32String(t)
		if err != nil {
			return nil, err
		}
		toAddresses = append(toAddresses, addr.AddressAsBech32String())
	}

	egldValue, err := decimal.NewFromString("1000000000000000000")
	if err != nil {
		return nil, err
	}

	args := blockchain.ArgsProxy{
		ProxyURL:            m.ProxyUrl,
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
		EntityType:          mxcore.Proxy,
	}
	ep, err := blockchain.NewProxy(args)
	if err != nil {
		return nil, err
	}
	w := interactors.NewWallet()

	bytePrivateKey, err := utils.ReadPrivateKey(pemPath)
	if err != nil {
		return nil, err
	}
	privateKey, err := w.LoadPrivateKeyFromPemData(bytePrivateKey)
	if err != nil {
		return nil, err
	}
	address, err := w.GetAddressFromPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	netConfigs, err := ep.GetNetworkConfig(context.Background())
	if err != nil {
		return nil, err
	}

	txBuilder, err := builders.NewTxBuilder(cryptoProvider.NewSigner())
	if err != nil {
		return nil, err
	}

	ti, err := interactors.NewTransactionInteractor(ep, txBuilder)
	if err != nil {
		return nil, err
	}

	transactionArguments, err := ep.GetDefaultTransactionArguments(context.Background(), address, netConfigs)
	if err != nil {
		return nil, err
	}

	var nonce = transactionArguments.Nonce

	var txs []*data.Transaction
	for _, toAddress := range toAddresses {
		transactionArguments.RcvAddr = toAddress
		transactionArguments.Value = amount.Mul(egldValue).String()

		transactionArguments.Data = []byte(dataStr)
		transactionArguments.GasLimit = netConfigs.MinGasLimit * 10

		transactionArguments.Version = 2
		transactionArguments.Options = 1
		transactionArguments.Nonce = nonce

		holder, _ := cryptoProvider.NewCryptoComponentsHolder(keyGen, privateKey)
		tx, err := ti.ApplySignatureAndGenerateTx(holder, transactionArguments)
		if err != nil {
			return nil, err
		}
		ti.AddTransaction(tx)
		txs = append(txs, tx)
		nonce++
	}

	hashes, err := ti.SendTransactions(context.Background(), txs)
	if err != nil {
		return nil, err
	}
	return hashes, nil
}

func (m MultiversxService) GetTrxStatus(hash string) (string, error) {
	args := blockchain.ArgsProxy{
		ProxyURL:            m.ProxyUrl,
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
		EntityType:          mxcore.Proxy,
	}
	ep, err := blockchain.NewProxy(args)
	if err != nil {
		return "", err
	}

	for {
		status, err := ep.GetTransactionStatus(context.Background(), hash)
		if status != "pending" && err == nil {
			if status == "invalid" {
				return "fail", nil
			}
			return status, nil
		}

		time.Sleep(1 * time.Second)
	}
}

func (m MultiversxService) GetAccount(address string) (float64, error) {
	args := blockchain.ArgsProxy{
		ProxyURL:            m.ProxyUrl,
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
		EntityType:          mxcore.Proxy,
	}
	ep, err := blockchain.NewProxy(args)
	if err != nil {
		return 0, err
	}

	networkConfig, err := ep.GetNetworkConfig(context.Background())
	if err != nil {
		return 0, err
	}

	addr, err := data.NewAddressFromBech32String(address)
	if err != nil {
		return 0, err
	}

	accountInfo, err := ep.GetAccount(context.Background(), addr)
	if err != nil {
		return 0, err
	}

	floatBalance, err := accountInfo.GetBalance(networkConfig.Denomination)
	if err != nil {
		return 0, err
	}
	floatBalance = floatBalance - 0.02
	return floatBalance, nil
}

func (m MultiversxService) CreateSwapTokensFixedInput(pemPath string, contract string, fromToken string, amount decimal.Decimal, toToken string, slippage decimal.Decimal, fromDecimals int, toDecimals int) (core.SwapResult, error) {
	toAddress, err := data.NewAddressFromBech32String(contract)
	if err != nil {
		return core.SwapResult{}, err
	}

	fromEgldValueStr := "1"
	for i := 0; i < fromDecimals; i++ {
		fromEgldValueStr += "0"
	}

	toEgldValueStr := "1"
	for i := 0; i < toDecimals; i++ {
		toEgldValueStr += "0"
	}

	fromEgldValue, err := decimal.NewFromString(fromEgldValueStr)
	if err != nil {
		return core.SwapResult{}, err
	}

	toEgldValue, err := decimal.NewFromString(toEgldValueStr)
	if err != nil {
		return core.SwapResult{}, err
	}
	amount = amount.Mul(fromEgldValue)
	slippage = slippage.Mul(toEgldValue)

	args := blockchain.ArgsProxy{
		ProxyURL:            m.ProxyUrl,
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
		EntityType:          mxcore.Proxy,
	}
	ep, err := blockchain.NewProxy(args)

	if err != nil {
		return core.SwapResult{}, err
	}
	w := interactors.NewWallet()

	bytePrivateKey, err := utils.ReadPrivateKey(pemPath)
	if err != nil {
		return core.SwapResult{}, err
	}
	privateKey, err := w.LoadPrivateKeyFromPemData(bytePrivateKey)
	if err != nil {
		return core.SwapResult{}, err
	}
	// Generate address from private key
	address, err := w.GetAddressFromPrivateKey(privateKey)
	if err != nil {
		return core.SwapResult{}, err
	}

	netConfigs, err := ep.GetNetworkConfig(context.Background())
	if err != nil {
		return core.SwapResult{}, err
	}

	transactionArguments, err := ep.GetDefaultTransactionArguments(context.Background(), address, netConfigs)
	if err != nil {
		return core.SwapResult{}, err
	}
	transactionArguments.RcvAddr = toAddress.AddressAsBech32String()
	transactionArguments.Value = "0"

	data := m.SwapTokensFixedInputData(fromToken, amount, toToken, slippage)
	transactionArguments.Data = []byte(data)
	transactionArguments.GasLimit = 50000000

	txBuilder, err := builders.NewTxBuilder(cryptoProvider.NewSigner())
	if err != nil {
		return core.SwapResult{}, err
	}

	ti, err := interactors.NewTransactionInteractor(ep, txBuilder)
	if err != nil {
		return core.SwapResult{}, err
	}
	holder, _ := cryptoProvider.NewCryptoComponentsHolder(keyGen, privateKey)
	tx, err := ti.ApplySignatureAndGenerateTx(holder, transactionArguments)
	if err != nil {
		return core.SwapResult{}, err
	}
	ti.AddTransaction(tx)

	hashes, err := ti.SendTransaction(context.Background(), tx)
	if err != nil {
		return core.SwapResult{}, err
	}

	status, err := m.GetContractStatus(hashes)
	if err != nil {
		return core.SwapResult{}, err
	}
	return core.SwapResult{
		Status: status,
		Hash:   hashes,
	}, nil
}

func (m MultiversxService) SwapTokensFixedInputData(fromToken string, amount decimal.Decimal, toToken string, slippage decimal.Decimal) string {
	fromTokenHex := hex.EncodeToString([]byte(fromToken))
	amountHex := fmt.Sprintf("%x", amount.BigInt())
	if len(amountHex)%2 != 0 {
		amountHex = "0" + amountHex
	}
	methodHex := hex.EncodeToString([]byte("swapTokensFixedInput"))
	toTokenHex := hex.EncodeToString([]byte(toToken))
	slippageHex := fmt.Sprintf("%x", slippage.BigInt())
	if len(slippageHex)%2 != 0 {
		slippageHex = "0" + slippageHex
	}

	finalStr := "ESDTTransfer@" + string(fromTokenHex) + "@" + string(amountHex) + "@" + string(methodHex) + "@" + string(toTokenHex) + "@" + string(slippageHex)
	return finalStr
}

func (m MultiversxService) GetContractStatus(hash string) (string, error) {
	args := blockchain.ArgsProxy{
		ProxyURL:            m.ProxyUrl,
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
		EntityType:          mxcore.Proxy,
	}
	ep, err := blockchain.NewProxy(args)
	if err != nil {
		return "fail", err
	}

	for {
		result, err := ep.GetTransactionInfoWithResults(context.Background(), hash)
		// fmt.Printf("%+v %+v\n", result, err)

		if err == nil {
			if result.Data.Transaction.Status != "pending" {
				if len(result.Data.Transaction.ScResults) > 0 {
					return "success", nil
				} else {
					return "fail", nil
				}
			}
		}

		time.Sleep(1 * time.Second)
	}
}
