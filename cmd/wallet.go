/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
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

		controller.CreateWallet(passwordFlag, amountFlag, mode)
	},
}

func init() {
	rootCmd.AddCommand(walletCmd)

	walletCmd.Flags().IntP("amount", "a", 1, "Number of wallets to create")
	walletCmd.Flags().StringP("password", "P", "Password123", "Default password for json wallet file")
}
