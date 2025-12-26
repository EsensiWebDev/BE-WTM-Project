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

// FormatCurrency formats a price value with currency symbol and proper formatting
// For IDR: uses dot (.) as thousands separator, no decimal places
// Format: "IDR 750.000"
func FormatCurrency(amount float64, currencyCode string, symbol string) string {
	normalizedCurrency := NormalizeCurrencyCode(currencyCode)
	decimalPlaces := GetDecimalPlaces(normalizedCurrency)

	// Round to appropriate decimal places
	var roundedAmount float64
	if decimalPlaces == 0 {
		roundedAmount = float64(int64(amount + 0.5)) // Round to nearest integer
	} else {
		// For currencies with decimals, round to specified decimal places
		multiplier := 1.0
		for i := 0; i < decimalPlaces; i++ {
			multiplier *= 10
		}
		roundedAmount = float64(int64(amount*multiplier+0.5)) / multiplier
	}

	// Format number with appropriate decimal places
	var amountStr string
	if decimalPlaces == 0 {
		amountStr = fmt.Sprintf("%.0f", roundedAmount)
	} else {
		amountStr = fmt.Sprintf("%."+fmt.Sprintf("%d", decimalPlaces)+"f", roundedAmount)
	}

	// Add thousands separators
	// For IDR: use dot (.) as thousands separator, no decimal separator needed
	// For others: use comma (,) as thousands separator, dot (.) as decimal separator
	var formattedAmount string
	if normalizedCurrency == "IDR" {
		// Indonesian format: dot as thousands separator, no decimals
		formattedAmount = formatWithSeparator(amountStr, ".", "")
	} else {
		// Default: comma as thousands separator, dot as decimal separator
		formattedAmount = formatWithSeparator(amountStr, ",", ".")
	}

	// Return formatted string with symbol prefix
	return fmt.Sprintf("%s %s", symbol, formattedAmount)
}

// formatWithSeparator adds thousands separator to a number string
// numStr: the number string (e.g., "750000" or "750000.50")
// thousandsSep: the thousands separator (e.g., "." or ",")
// decimalSep: the decimal separator (e.g., "." or ","), empty string if no decimals
func formatWithSeparator(numStr string, thousandsSep string, decimalSep string) string {
	// Split integer and decimal parts using the decimal separator
	var integerPart, decimalPart string
	if decimalSep != "" {
		parts := strings.Split(numStr, decimalSep)
		integerPart = parts[0]
		if len(parts) > 1 {
			decimalPart = decimalSep + parts[1]
		}
	} else {
		// No decimal separator, entire string is integer part
		integerPart = numStr
		decimalPart = ""
	}

	// Add thousands separators from right to left
	var result strings.Builder
	count := 0
	for i := len(integerPart) - 1; i >= 0; i-- {
		if count > 0 && count%3 == 0 {
			result.WriteString(thousandsSep)
		}
		result.WriteByte(integerPart[i])
		count++
	}

	// Reverse the string
	reversed := result.String()
	formattedInteger := ""
	for i := len(reversed) - 1; i >= 0; i-- {
		formattedInteger += string(reversed[i])
	}

	return formattedInteger + decimalPart
}
