package controller

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"github.com/xm1k3/mxga/pkg/core"
	"github.com/xm1k3/mxga/pkg/core/services"
	"github.com/xm1k3/mxga/pkg/utils"
)

type Controller struct {
	Service core.MultiversxNetService
}

func GetController(mode string) Controller {
	multiversxService := services.MultiversxService{
		ProxyUrl: utils.GetProxyUrl(mode),
	}
	controller := Controller{
		Service: multiversxService,
	}
	return controller
}

func CreateWallet(password string, qty int, mode string) {
	controller := GetController(mode)

	savePath := filepath.Dir(viper.ConfigFileUsed())
	savePath = path.Join(savePath, "wallets")
	for i := 0; i < qty; i++ {
		wallet, err := controller.Service.CreateWallet(viper.ConfigFileUsed(), password)
		if err != nil {
			log.Fatal(err)
		}
		reportData := []byte("Wallet address: " + wallet.Address + "\nSecret words:" + wallet.Mnemonic + "\nPem file path:" + wallet.PemPath + "\nJson file path:" + wallet.JsonPath + "\nJson password:" + wallet.JsonPassword + "\n\n")

		f, err := os.OpenFile(path.Join(savePath, "report.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f.Write(reportData); err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Wallet created:" + wallet.Address)
	}

	fmt.Println("Report created here: ", path.Join(savePath, "report.txt"))
}

func SendTransactions(pemPath string, to []string, amount decimal.Decimal, data string, mode string) {
	controller := GetController(mode)
	hashes, err := controller.Service.SendTransactions(pemPath, to, amount, data)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(to))
	for i := 0; i < len(hashes); i++ {
		go func(i int) {
			defer wg.Done()
			status, err := controller.Service.GetTrxStatus(hashes[i])
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("[", status, "] Hash: ", hashes[i])
		}(i)
	}
	wg.Wait()
}

func Retrieve(walletsAddr []string, walletsPemPath []string, mainAddress string, amount decimal.Decimal, datastr string, mode string, all bool) {
	controller := GetController(mode)
	var hashes []string
	for i, wallet := range walletsPemPath {
		var mainStr []string

		// retrieve all account data
		if all {
			amountStr, err := controller.Service.GetAccount(walletsAddr[i])
			if err != nil {
				log.Fatal(err)
			}
			amount, err = decimal.NewFromString(amountStr)
			if err != nil {
				log.Fatal(err)
			}
		}

		mainStr = append(mainStr, mainAddress)
		hash, err := controller.Service.SendTransactions(wallet, mainStr, amount, datastr)
		if err != nil {
			log.Fatal(err)
		}
		hashes = append(hashes, hash...)
	}

	var wg sync.WaitGroup
	wg.Add(len(hashes))
	for i := 0; i < len(hashes); i++ {
		go func(i int) {
			defer wg.Done()
			status, err := controller.Service.GetTrxStatus(hashes[i])
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("[", status, "] Hash: ", hashes[i])
		}(i)
	}
	wg.Wait()
}
