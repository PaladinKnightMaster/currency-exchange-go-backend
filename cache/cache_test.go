package cache

import (
	"testing"
	"time"

	"github.com/paladinknightmaster/currency-exchange-go-backend/models"
)

func TestCache_SaveAndGetRates(t *testing.T) {
	ClearCache()

	mockRates := &models.ExchangeRates{
		BaseCode: "USD",
		Rates: map[string]float64{
			"EUR": 0.85,
		},
	}

	SaveRatesToCache(mockRates)

	rates, found := GetRatesFromCache()
	if !found {
		t.Error("Expected to find rates in cache")
	}

	if rates.Rates["EUR"] != 0.85 {
		t.Errorf("Expected rate 0.85 for EUR, got %f", rates.Rates["EUR"])
	}
}

func TestCache_Expiration(t *testing.T) {
	ClearCache()
	SetCacheExpiration(1 * time.Second) // Set expiration to 1 second

	mockRates := &models.ExchangeRates{
		BaseCode: "USD",
		Rates:    map[string]float64{"EUR": 0.85},
	}

	SaveRatesToCache(mockRates)

	time.Sleep(2 * time.Second) // Wait for cache to expire

	_, found := GetRatesFromCache()
	if found {
		t.Error("Expected cache to be expired and not find rates")
	}

	// Reset cache expiration to default
	SetCacheExpiration(10 * time.Minute)
}
