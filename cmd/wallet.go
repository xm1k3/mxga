/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xm1k3/mxga/pkg/core/controller"
)

// walletCmd represents the wallet command
var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Create new wallets",
	Long:  `Create new wallets`,
	Run: func(cmd *cobra.Command, args []string) {
		amountFlag, _ := cmd.Flags().GetInt("amount")
		passwordFlag, _ := cmd.Flags().GetString("password")
		mode, _ := rootCmd.PersistentFlags().GetString("mode")

		configFileUsed := viper.ConfigFileUsed()
		if configFileUsed == "" {
			wd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			configFileUsed = wd
		}

		controller.CreateWallet(configFileUsed, passwordFlag, amountFlag, mode)
	},
}

func init() {
	rootCmd.AddCommand(walletCmd)

	walletCmd.Flags().IntP("amount", "a", 1, "Number of wallets to create")
	walletCmd.Flags().StringP("password", "p", "Password123", "Default password for json wallet file")
}
