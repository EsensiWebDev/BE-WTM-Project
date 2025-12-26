package seed

import (
	"log"
	"wtm-backend/internal/infrastructure/database/model"
)

func (s *Seed) SeedCurrency() {
	var countCurrency int64
	s.db.Model(&model.Currency{}).Count(&countCurrency)
	if countCurrency == 0 {
		currencies := []model.Currency{
			{Code: "IDR", Name: "Indonesian Rupiah", Symbol: "IDR", IsActive: true},
			{Code: "USD", Name: "US Dollar", Symbol: "$", IsActive: true},
			{Code: "EUR", Name: "Euro", Symbol: "€", IsActive: true},
			{Code: "GBP", Name: "British Pound", Symbol: "£", IsActive: true},
			{Code: "JPY", Name: "Japanese Yen", Symbol: "¥", IsActive: true},
			{Code: "KRW", Name: "South Korean Won", Symbol: "₩", IsActive: true},
			{Code: "SGD", Name: "Singapore Dollar", Symbol: "S$", IsActive: true},
			{Code: "MYR", Name: "Malaysian Ringgit", Symbol: "RM", IsActive: true},
			{Code: "THB", Name: "Thai Baht", Symbol: "฿", IsActive: true},
			{Code: "CNY", Name: "Chinese Yuan", Symbol: "¥", IsActive: true},
			{Code: "AUD", Name: "Australian Dollar", Symbol: "A$", IsActive: true},
			{Code: "CAD", Name: "Canadian Dollar", Symbol: "C$", IsActive: true},
			{Code: "CHF", Name: "Swiss Franc", Symbol: "CHF", IsActive: true},
			{Code: "HKD", Name: "Hong Kong Dollar", Symbol: "HK$", IsActive: true},
			{Code: "NZD", Name: "New Zealand Dollar", Symbol: "NZ$", IsActive: true},
		}

		if err := s.db.Create(&currencies).Error; err != nil {
			log.Fatalf("Failed to seed currencies: %s", err.Error())
		}
		log.Println("Seeding currencies completed")
	}
}
