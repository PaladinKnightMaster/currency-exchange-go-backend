package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/paladinknightmaster/currency-exchange-go-backend/models"
)

func TestFetchRates_Success(t *testing.T) {
	// Mock external API response
	mockRates := &models.ExchangeRates{
		BaseCode: "USD",
		Rates: map[string]float64{
			"EUR": 0.85,
		},
	}

	mockResponseBody, _ := json.Marshal(mockRates)

	Client = &MockClient{
		ResponseBody: string(mockResponseBody),
		StatusCode:   http.StatusOK,
	}

	rates, err := FetchRates()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rates.Rates["EUR"] != 0.85 {
		t.Errorf("Expected rate 0.85 for EUR, got %f", rates.Rates["EUR"])
	}
}

func TestFetchRates_Error(t *testing.T) {
	Client = &MockClient{
		Err: errors.New("network error"),
	}

	_, err := FetchRates()
	if err == nil {
		t.Error("Expected an error, got nil")
	}
}

func TestFetchRates_InvalidJSON(t *testing.T) {
	Client = &MockClient{
		ResponseBody: "invalid json",
		StatusCode:   http.StatusOK,
	}

	_, err := FetchRates()
	if err == nil {
		t.Error("Expected an error due to invalid JSON, got nil")
	}
}

func TestFetchRates_Non200StatusCode(t *testing.T) {
	Client = &MockClient{
		ResponseBody: "",
		StatusCode:   http.StatusInternalServerError,
	}

	_, err := FetchRates()
	if err == nil {
		t.Error("Expected an error due to non-200 status code, got nil")
	}
}

func TestFetchRates_NoAPIKey(t *testing.T) {
	// Unset the API key
	os.Unsetenv("EXCHANGE_RATE_API_KEY")

	_, err := FetchRates()
	if err == nil {
		t.Error("Expected an error due to missing API key, got nil")
	}
}
