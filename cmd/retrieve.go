/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"log"

	"github.com/multiversx/mx-sdk-go/interactors"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xm1k3/mxga/pkg/core"
	"github.com/xm1k3/mxga/pkg/core/controller"
	"github.com/xm1k3/mxga/pkg/utils"
)

// retrieveCmd represents the retrieve command
var retrieveCmd = &cobra.Command{
	Use:   "retrieve",
	Short: "retrieve back money from all wallets",
	Long:  `retrieve back money from all wallets`,
	Run: func(cmd *cobra.Command, args []string) {
		valueFlag, _ := cmd.Flags().GetFloat32("value")
		dataFlag, _ := cmd.Flags().GetString("data")
		allFlag, _ := cmd.Flags().GetBool("all")
		mode := viper.GetString("mode")

		wallet := viper.GetStringMap("wallet")
		walletJson, err := json.Marshal(wallet)
		if err != nil {
			log.Fatal(err)
		}

		var confs core.WalletConfig
		if err := json.Unmarshal(walletJson, &confs); err != nil {
			log.Fatal(err)
		}

		amount := decimal.NewFromFloat32(valueFlag)
		mainWalletPemPath := confs.Path + confs.Main + confs.Ext

		// pem addresses
		var walletsPemPath []string
		// single wallet address
		var walletsAddr []string
		for _, elem := range confs.Other {
			walletsPemPath = append(walletsPemPath, confs.Path+elem+confs.Ext)
			walletsAddr = append(walletsAddr, elem)
		}

		w := interactors.NewWallet()
		bytePrivateKey, err := utils.ReadPrivateKey(mainWalletPemPath)
		if err != nil {
			log.Fatal(err)
		}
		privateKey, err := w.LoadPrivateKeyFromPemData(bytePrivateKey)
		if err != nil {
			log.Fatal(err)
		}
		mainAddress, err := w.GetAddressFromPrivateKey(privateKey)
		if err != nil {
			log.Fatal(err)
		}

		controller.Retrieve(walletsAddr, walletsPemPath, mainAddress.AddressAsBech32String(), amount, dataFlag, mode, allFlag)

	},
}

func init() {
	rootCmd.AddCommand(retrieveCmd)

	retrieveCmd.Flags().Float32P("value", "v", 0.5, "value")
	retrieveCmd.Flags().StringP("data", "d", "", "data")
	retrieveCmd.Flags().BoolP("all", "a", false, "retrieve all money from all wallets")

}
