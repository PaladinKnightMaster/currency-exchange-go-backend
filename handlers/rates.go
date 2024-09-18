package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/paladinknightmaster/currency-exchange-go-backend/cache"
	"github.com/paladinknightmaster/currency-exchange-go-backend/utils"
)

func GetRates(w http.ResponseWriter, r *http.Request) {
	rates, found := cache.GetRatesFromCache()
	if !found {
		// Fetch rates from external API
		fetchedRates, err := utils.FetchRatesFunc()
		if err != nil {
			http.Error(w, "Error fetching rates", http.StatusInternalServerError)
			return
		}
		rates = fetchedRates
		cache.SaveRatesToCache(rates)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rates)
}
