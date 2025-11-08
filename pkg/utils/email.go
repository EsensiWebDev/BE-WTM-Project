package utils

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

func ParseTemplate(tmpl string, data interface{}) (string, error) {
	funcMap := template.FuncMap{
		"add1": func(i int) int { return i + 1 },
	}
	t, err := template.New("email").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	return buf.String(), err
}

func HumanizeDuration(d time.Duration) string {
	// Konversi ke detik total
	seconds := int(d.Seconds())

	if seconds < 60 {
		return fmt.Sprintf("%d seconds", seconds)
	}

	minutes := seconds / 60
	if minutes < 60 {
		return fmt.Sprintf("%d minutes", minutes)
	}

	hours := minutes / 60
	if hours < 24 {
		remainingMinutes := minutes % 60
		if remainingMinutes == 0 {
			return fmt.Sprintf("%d hours", hours)
		}
		return fmt.Sprintf("%d hours %d minutes", hours, remainingMinutes)
	}

	days := hours / 24
	remainingHours := hours % 24
	if remainingHours == 0 {
		return fmt.Sprintf("%d days", days)
	}
	return fmt.Sprintf("%d days %d hours", days, remainingHours)
}
