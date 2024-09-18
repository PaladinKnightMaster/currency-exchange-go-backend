package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/paladinknightmaster/currency-exchange-go-backend/models"
)

var Client HTTPClient = &http.Client{Timeout: 10 * time.Second}

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

var FetchRatesFunc = FetchRates

func FetchRates() (*models.ExchangeRates, error) {
	apiKey := os.Getenv("EXCHANGE_RATE_API_KEY")
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/USD", apiKey)

	resp, err := Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	var rates models.ExchangeRates
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		return nil, err
	}

	return &rates, nil
}
