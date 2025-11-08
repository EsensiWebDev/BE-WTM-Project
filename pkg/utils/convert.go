package utils

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	"wtm-backend/pkg/constant"
)

func StringToUint(s string) (uint, error) {
	if s == "" {
		return 0, fmt.Errorf("empty string cannot be converted to uint")
	}
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		log.Printf("StringToUint error: failed to convert '%s': %s", s, err.Error())
		return 0, err
	}

	if val == 0 {
		return 0, fmt.Errorf("the value cannot be zero")
	}

	return uint(val), nil
}

// ParseRFC3339ToUTC parses RFC3339 string with timezone offset
// and converts to UTC for safe DB insert into `timestamp with time zone`
func ParseRFC3339ToUTC(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("empty datetime string")
	}
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid RFC3339 format: %s", err.Error())
	}
	return t.UTC(), nil
}

func DaysInMonth(monthTime time.Time) (int, error) {
	monthStart := monthTime
	nextMonth := monthStart.AddDate(0, 1, 0)

	// Selisih hari antara awal bulan sekarang dan awal bulan depan
	duration := nextMonth.Sub(monthStart)
	return int(duration.Hours() / 24), nil
}

// Regex: 4 digit tahun + tanda "-" + 2 digit bulan (01-12)
var monthRegex = regexp.MustCompile(`^\d{4}-(0[1-9]|1[0-2])$`)

// IsValidMonth memverifikasi format bulan sebagai string "YYYY-MM"
func IsValidMonth(month string) bool {
	return monthRegex.MatchString(month)
}

func ParseHourString(s string) (*time.Time, error) {
	// Gabungkan dengan tanggal dummy agar zona waktu bisa dihitung akurat
	full := fmt.Sprintf("2000-01-01T%s:00", s) // contoh: "2000-01-01T14:00:00"
	t, err := time.ParseInLocation("2006-01-02T15:04:05", full, constant.AsiaJakarta)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func CapitalizeWords(s string) string {
	words := strings.Fields(s)
	for i, w := range words {
		words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
	}
	return strings.Join(words, " ")
}
