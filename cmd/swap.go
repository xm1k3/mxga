/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xm1k3/mxga/pkg/core"
	"github.com/xm1k3/mxga/pkg/core/controller"
)

// swapCmd represents the swap command
var swapCmd = &cobra.Command{
	Use:   "swap",
	Short: "Manual swap",
	Long:  `Do a swap transactions manually`,
	Run: func(cmd *cobra.Command, args []string) {
		contract, _ := cmd.Flags().GetString("contract")
		fromToken, _ := cmd.Flags().GetString("from")
		toToken, _ := cmd.Flags().GetString("to")
		amount, _ := cmd.Flags().GetFloat32("amount")
		slippage, _ := cmd.Flags().GetFloat32("slippage")
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

		mainWalletPemPath := filepath.Join(confs.Path, confs.Main+confs.Ext)

		if mode != "manual" {
			contract = viper.GetString(mode + ".contract")
			fromToken = viper.GetString(mode + ".swap.from")
			toToken = viper.GetString(mode + ".swap.to")
			slippage = float32(viper.GetFloat64(mode + ".swap.slippage"))
			amount = float32(viper.GetFloat64(mode + ".swap.amount"))
		}
		if contract == "" {
			log.Fatal("Insert a contract address")
		}
		if allFlag {
			amount = float32(controller.GetAccountTokenPrice(confs.Address, fromToken, mode).InexactFloat64())
		}

		controller.CreateSwapTokensFixedInput(mainWalletPemPath, contract, fromToken, amount, toToken, slippage, mode)

	},
}

func init() {
	rootCmd.AddCommand(swapCmd)
	swapCmd.Flags().StringP("contract", "c", "", "contract where to swap")
	swapCmd.Flags().StringP("from", "f", "", "from")
	swapCmd.Flags().StringP("to", "t", "", "to")
	swapCmd.Flags().Float32P("amount", "a", 0.5, "amount")
	swapCmd.Flags().Float32P("slippage", "s", 1, "slippage")
	swapCmd.Flags().BoolP("all", "", false, "buy all")
}
