package currency

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Valid ISO 4217 currency codes (common ones)
var validCurrencyCodes = map[string]bool{
	"IDR": true, "USD": true, "EUR": true, "GBP": true, "JPY": true,
	"KRW": true, "SGD": true, "MYR": true, "THB": true, "CNY": true,
	"AUD": true, "CAD": true, "CHF": true, "HKD": true, "NZD": true,
	"INR": true, "PHP": true, "VND": true, "BRL": true, "MXN": true,
	"ZAR": true, "RUB": true, "TRY": true, "SEK": true, "NOK": true,
	"DKK": true, "PLN": true, "CZK": true, "HUF": true, "ILS": true,
	"AED": true, "SAR": true, "QAR": true, "KWD": true, "BHD": true,
}

// ValidateCurrencyCode validates if a currency code is a valid ISO 4217 code
func ValidateCurrencyCode(code string) bool {
	if code == "" {
		return false
	}
	normalized := NormalizeCurrencyCode(code)
	return validCurrencyCodes[normalized]
}

// NormalizeCurrencyCode normalizes currency code to uppercase
func NormalizeCurrencyCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

// ValidatePricesHasIDR validates that prices map contains IDR (mandatory currency)
func ValidatePricesHasIDR(prices map[string]float64) bool {
	if prices == nil {
		return false
	}
	_, exists := prices["IDR"]
	return exists
}

// GetPriceForCurrency gets price for a specific currency, with fallback to IDR
func GetPriceForCurrency(prices map[string]float64, currency string) (float64, string, error) {
	if prices == nil || len(prices) == 0 {
		return 0, "", fmt.Errorf("prices map is empty")
	}

	normalizedCurrency := NormalizeCurrencyCode(currency)

	// Try requested currency first
	if price, exists := prices[normalizedCurrency]; exists {
		return price, normalizedCurrency, nil
	}

	// Fallback to IDR (mandatory currency)
	if price, exists := prices["IDR"]; exists {
		return price, "IDR", fmt.Errorf("currency %s not available, using IDR", currency)
	}

	return 0, "", fmt.Errorf("no price available for currency %s and IDR fallback not found", currency)
}

// ValidatePrices validates that prices map is valid
// - Must have at least IDR
// - All prices must be positive
// - All currency codes must be valid
func ValidatePrices(prices map[string]float64) error {
	if prices == nil || len(prices) == 0 {
		return fmt.Errorf("prices map cannot be empty")
	}

	// Check IDR exists
	if !ValidatePricesHasIDR(prices) {
		return fmt.Errorf("prices must contain IDR (mandatory currency)")
	}

	// Validate all currency codes and prices
	for code, price := range prices {
		normalizedCode := NormalizeCurrencyCode(code)
		if !ValidateCurrencyCode(normalizedCode) {
			return fmt.Errorf("invalid currency code: %s", code)
		}
		if price < 0 {
			return fmt.Errorf("price cannot be negative for currency %s", code)
		}
	}

	return nil
}

// PricesToJSON converts prices map to JSON bytes
func PricesToJSON(prices map[string]float64) ([]byte, error) {
	return json.Marshal(prices)
}

// JSONToPrices converts JSON bytes to prices map
func JSONToPrices(data []byte) (map[string]float64, error) {
	var prices map[string]float64
	if err := json.Unmarshal(data, &prices); err != nil {
		return nil, err
	}
	return prices, nil
}

// GetDecimalPlaces returns the number of decimal places for a currency
// Most currencies use 2 decimal places, but some use 0 (JPY, KRW, IDR)
func GetDecimalPlaces(currency string) int {
	normalized := NormalizeCurrencyCode(currency)
	zeroDecimalCurrencies := map[string]bool{
		"JPY": true, "KRW": true, "IDR": true, "VND": true,
	}
	if zeroDecimalCurrencies[normalized] {
		return 0
	}
	return 2
}
