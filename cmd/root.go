/*
Copyright © 2023 xm1k3
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mxga",
	Short: "MultiversX Golang Api, tool for interact with API & Blockchain",
	Long:  `MultiversX Golang Api, tool for interact with API & Blockchain`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("mode", "M", "mainnet", "multiversx mode (mainnet, testnet, devnet)")
}
