package services

import (
	"context"
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
