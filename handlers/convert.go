package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/paladinknightmaster/currency-exchange-go-backend/cache"
)

func ConvertCurrency(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	from := query.Get("from")
	to := query.Get("to")
	amountStr := query.Get("amount")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	rates, found := cache.GetRatesFromCache()
	if !found {
		http.Error(w, "Rates not available", http.StatusServiceUnavailable)
		return
	}

	fromRate, okFrom := rates.Rates[from]
	toRate, okTo := rates.Rates[to]

	if !okFrom || !okTo {
		http.Error(w, "Invalid currency code", http.StatusBadRequest)
		return
	}

	convertedAmount := (amount / fromRate) * toRate
	// Round the converted amount to 2 decimal places
	convertedAmount = math.Round(convertedAmount*100) / 100

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"converted_amount": convertedAmount})
}
