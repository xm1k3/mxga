package services

import (
	"os"
	"path"
	"path/filepath"

	"github.com/multiversx/mx-sdk-go/interactors"
	"github.com/xm1k3/mxga/pkg/core"
)

type MultiversxService struct {
	ProxyUrl string
}

// Function for create a single wallet, returning a struct
func (e MultiversxService) CreateWallet(folderPath string, password string) (core.Wallet, error) {
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
