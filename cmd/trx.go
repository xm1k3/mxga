/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"log"

	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xm1k3/mxga/pkg/core"
	"github.com/xm1k3/mxga/pkg/core/controller"
)

// trxCmd represents the trx command
var trxCmd = &cobra.Command{
	Use:   "trx",
	Short: "Send multiple transactions",
	Long:  `Send multiple transactions`,
	Run: func(cmd *cobra.Command, args []string) {
		valueFlag, _ := cmd.Flags().GetFloat32("value")
		dataFlag, _ := cmd.Flags().GetString("data")
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

		controller.SendTransactions(mainWalletPemPath, confs.Other, amount, dataFlag, mode)
	},
}

func init() {
	rootCmd.AddCommand(trxCmd)

	trxCmd.Flags().Float32P("value", "v", 0.1, "value")
	trxCmd.Flags().StringP("data", "d", "", "data")
}
