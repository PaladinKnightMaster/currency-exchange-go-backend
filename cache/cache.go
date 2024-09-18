package cache

import (
	"time"

	"github.com/paladinknightmaster/currency-exchange-go-backend/models"
	"github.com/patrickmn/go-cache"
)

var c = cache.New(10*time.Minute, 15*time.Minute)

func GetRatesFromCache() (*models.ExchangeRates, bool) {
	rates, found := c.Get("rates")
	if found {
		return rates.(*models.ExchangeRates), true
	}
	return nil, false
}

func SaveRatesToCache(rates *models.ExchangeRates) {
	c.Set("rates", rates, cache.DefaultExpiration)
}

func ClearCache() {
	c.Flush()
}

func SetCacheExpiration(d time.Duration) {
	c = cache.New(d, 15*time.Minute)
}
