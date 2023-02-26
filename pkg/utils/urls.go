package utils

import (
	"log"

	"github.com/xm1k3/mxga/pkg/core"
)

func GetProxyUrl(mode string) string {
	switch mode {
	case "testnet":
		return "https://testnet-gateway.multiversx.com"
	case "devnet":
		return "https://devnet-gateway.multiversx.com"
	case "mainnet":
		return "https://gateway.multiversx.com"
	default:
		log.Fatal(core.ErrInvalidMode)
		return ""
	}
}

func GetApiUrl(mode string) string {
	switch mode {
	case "testnet":
		return "https://testnet-api.multiversx.com"
	case "devnet":
		return "https://devnet-api.multiversx.com"
	case "mainnet":
		return "https://api.multiversx.com"
	default:
		log.Fatal(core.ErrInvalidMode)
		return ""
	}
}
