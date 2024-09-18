package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/paladinknightmaster/currency-exchange-go-backend/cache"
	"github.com/paladinknightmaster/currency-exchange-go-backend/models"
	"github.com/paladinknightmaster/currency-exchange-go-backend/utils"
)

func TestGetRates_Success(t *testing.T) {
	// Clear cache before test
	cache.ClearCache()

	// Mock external API response
	mockRates := &models.ExchangeRates{
		BaseCode: "USD",
		Rates: map[string]float64{
			"EUR": 0.85,
			"GBP": 0.75,
			"CAD": 1.25,
		},
	}

	mockResponseBody, _ := json.Marshal(mockRates)

	utils.Client = &utils.MockClient{
		ResponseBody: string(mockResponseBody),
		StatusCode:   http.StatusOK,
	}

	req, err := http.NewRequest("GET", "/rates", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetRates)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var ratesResponse models.ExchangeRates
	if err := json.NewDecoder(rr.Body).Decode(&ratesResponse); err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	if ratesResponse.BaseCode != "USD" {
		t.Errorf("Expected BaseCode USD, got %s", ratesResponse.BaseCode)
	}

	// Verify rates
	expectedRates := mockRates.Rates
	for currency, rate := range expectedRates {
		if ratesResponse.Rates[currency] != rate {
			t.Errorf("Expected rate %f for %s, got %f", rate, currency, ratesResponse.Rates[currency])
		}
	}
}

func TestGetRates_FetchError(t *testing.T) {
	// Clear the cache to force fetching rates
	cache.ClearCache()

	// Mock FetchRatesFunc to return an error
	utils.FetchRatesFunc = func() (*models.ExchangeRates, error) {
		return nil, errors.New("failed to fetch rates")
	}
	// Reset FetchRatesFunc after the test
	defer func() {
		utils.FetchRatesFunc = utils.FetchRates
	}()

	req, err := http.NewRequest("GET", "/rates", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetRates)
	handler.ServeHTTP(rr, req)

	// Check that the status code is 500 Internal Server Error
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	// Check that the response body contains the error message
	expectedMessage := "Error fetching rates\n"
	if rr.Body.String() != expectedMessage {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedMessage)
	}
}

func TestConvertCurrency_Success(t *testing.T) {
	// Set up cache with mock rates
	cache.ClearCache()
	mockRates := &models.ExchangeRates{
		BaseCode: "USD",
		Rates: map[string]float64{
			"USD": 1.0,
			"EUR": 0.85,
			"GBP": 0.75,
		},
	}
	cache.SaveRatesToCache(mockRates)

	req, err := http.NewRequest("GET", "/convert?from=USD&to=EUR&amount=100", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ConvertCurrency)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var result map[string]float64
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	expectedAmount := 85.0 // 100 USD * 0.85 EUR/USD
	if result["converted_amount"] != expectedAmount {
		t.Errorf("Expected converted amount %f, got %f", expectedAmount, result["converted_amount"])
	}
}

func TestConvertCurrency_MissingParameters(t *testing.T) {
	req, err := http.NewRequest("GET", "/convert?from=USD&amount=100", nil) // Missing 'to' parameter
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ConvertCurrency)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestConvertCurrency_TableDriven(t *testing.T) {
	// Set up cache with mock rates
	cache.ClearCache()
	mockRates := &models.ExchangeRates{
		BaseCode: "USD",
		Rates: map[string]float64{
			"USD": 1.0,
			"EUR": 0.85,
			"GBP": 0.75,
		},
	}
	cache.SaveRatesToCache(mockRates)

	tests := []struct {
		from          string
		to            string
		amount        string
		expectedCode  int
		expectedValue float64
	}{
		{"USD", "EUR", "100", http.StatusOK, 85.0},
		{"GBP", "EUR", "100", http.StatusOK, 113.33},
		{"USD", "INVALID", "100", http.StatusBadRequest, 0},
		{"USD", "EUR", "-50", http.StatusBadRequest, 0},
	}

	for _, tt := range tests {
		req, err := http.NewRequest("GET", "/convert", nil)
		if err != nil {
			t.Fatal(err)
		}

		q := req.URL.Query()
		q.Add("from", tt.from)
		q.Add("to", tt.to)
		q.Add("amount", tt.amount)
		req.URL.RawQuery = q.Encode()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ConvertCurrency)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != tt.expectedCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, tt.expectedCode)
		}

		if tt.expectedCode == http.StatusOK {
			var result map[string]float64
			if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
				t.Errorf("Error decoding response: %v", err)
			}

			// Compare the rounded values
			if result["converted_amount"] != tt.expectedValue {
				t.Errorf("Expected converted amount %.2f, got %.2f", tt.expectedValue, result["converted_amount"])
			}
		}
	}
}
