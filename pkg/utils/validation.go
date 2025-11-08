package utils

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"strings"
	"unicode"
)

func NotEmptyAfterTrim(fieldName string) validation.Rule {
	return validation.By(func(value interface{}) error {
		s, _ := value.(string)
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("%s cannot be blank or whitespace", fieldName)
		}
		return nil
	})
}

func ParseValidationErrors(err error) map[string]string {
	var errs validation.Errors
	if errors.As(err, &errs) {
		errorMap := make(map[string]string)
		for field, valErr := range errs {
			errorMap[field] = valErr.Error()
		}
		return errorMap
	}
	return nil
}

// EscapeAndNormalizeSearch membersihkan input pencarian agar aman dan konsisten
func EscapeAndNormalizeSearch(input string) string {
	// Trim spasi di awal dan akhir
	input = strings.TrimSpace(input)

	// Normalisasi ke lowercase (opsional, tergantung kebutuhan search case-insensitive)
	input = strings.ToLower(input)

	// Escape karakter backslash terlebih dahulu
	input = strings.ReplaceAll(input, `\`, `\\`)

	// Escape wildcard SQL
	input = strings.ReplaceAll(input, `%`, `\%`)
	input = strings.ReplaceAll(input, `_`, `\_`)

	// (Opsional) Hapus karakter non-printable
	input = strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, input)

	return input
}
