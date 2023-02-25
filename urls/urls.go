package urls

func GetProxyUrl(mode string) string {
	switch mode {
	case "testnet":
		return "https://testnet-gateway.elrond.com"
	case "devnet":
		return "https://devnet-gateway.elrond.com"
	case "mainnet":
		return "https://gateway.elrond.com"
	default:
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
		return ""
	}
}
