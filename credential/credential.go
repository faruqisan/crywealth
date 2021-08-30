package credential

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type (
	Binance struct {
		APIKey    string
		APISecret string
	}
)

var (
	ErrorCredentialsMissing = errors.New("missing binance credentials")
)

// LoadENVFromFile will read your env file(s) and load them into ENV for this process.
func LoadENVFromFile() error {
	return godotenv.Load()
}

func BinanceSecret() (Binance, error) {
	apiKey := os.Getenv("BINANCE_API_KEY")
	apiSecret := os.Getenv("BINANCE_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		return Binance{}, ErrorCredentialsMissing
	}

	return Binance{
		APIKey:    apiKey,
		APISecret: apiSecret,
	}, nil
}
