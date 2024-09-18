package models

type ExchangeRates struct {
	BaseCode string             `json:"base_code"`
	Rates    map[string]float64 `json:"conversion_rates"`
}
