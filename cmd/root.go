/*
Copyright © 2023 xm1k3
*/
package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var cfgFile string

type Config struct {
	Mode string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "mxga",
	Long: `
███╗   ███╗██╗  ██╗ ██████╗  █████╗ 
████╗ ████║╚██╗██╔╝██╔════╝ ██╔══██╗
██╔████╔██║ ╚███╔╝ ██║  ███╗███████║
██║╚██╔╝██║ ██╔██╗ ██║   ██║██╔══██║
██║ ╚═╝ ██║██╔╝ ██╗╚██████╔╝██║  ██║
╚═╝     ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝

MultiversX Golang Api, tool for interact with API & Blockchain
By xm1k3
`,
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
	cobra.OnInitialize(initConfig)

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().StringP("mode", "M", "mainnet", "multiversx mode (mainnet, testnet, devnet)")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/mxga/mxga.yaml)")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		homeDir, err := os.UserHomeDir()
		cobra.CheckErr(err)

		configDirPath := filepath.Join(homeDir, ".config")

		configDirPath = filepath.Join(configDirPath, "mxga")
		if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
			if err := os.MkdirAll(configDirPath, 0755); err != nil {
				cobra.CheckErr(err)
			}
		}

		configFilePath := filepath.Join(configDirPath, "mxga.yaml")
		if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
			config := Config{
				Mode: "mainnet",
			}
			file, err := os.Create(configFilePath)
			if err != nil {
				cobra.CheckErr(err)
			}
			defer file.Close()

			encoder := yaml.NewEncoder(file)
			if err := encoder.Encode(&config); err != nil {
				cobra.CheckErr(err)
			}
		}

		viper.AddConfigPath(configDirPath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("mxga")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
