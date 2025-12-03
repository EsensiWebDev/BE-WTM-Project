package hotel_usecase

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/pkg/logger"
)

func (hu *HotelUsecase) UploadHotel(ctx context.Context, req *hoteldto.UploadHotelRequest) error {
	// 2️⃣ Buka file CSV
	file, err := req.File.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';' // pakai ; sebagai delimiter utama
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read csv: %w", err)
	}

	// 3️⃣ Iterasi record (skip header & komentar)
	for i, row := range records {
		if len(row) == 0 || strings.HasPrefix(row[0], "#") {
			continue
		}
		if i == 0 {
			// header row
			continue
		}

		var rating int
		rating = 0
		if _, err := fmt.Sscanf(row[6], "%d", &rating); err != nil {
			logger.Error(ctx, "failed to parse rating: %w", err)
		}

		// 4️⃣ Map ke struct CreateHotelRequest
		var nearbyPlaces []hoteldto.NearbyPlace
		var socialMedias []hoteldto.SocialMedia
		var jsonSocialMedias, jsonNearbyPlaces []byte
		for _, place := range strings.Split(row[7], "|") {
			parts := strings.Split(place, ",")
			if len(parts) != 2 {
				break
			}

			dist, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err != nil {
				break
			}

			nearbyPlaces = append(nearbyPlaces, hoteldto.NearbyPlace{
				Name:     place,
				Distance: dist,
			})
		}

		jsonNearbyPlaces, err = json.Marshal(nearbyPlaces)
		if err != nil {
			logger.Error(ctx, "failed to marshal nearby places: %w", err)
		}

		if IsValidLinkWithKeyword(row[9], "tiktok") {
			socialMedias = append(socialMedias, hoteldto.SocialMedia{
				Platform: "Tiktok",
				Link:     row[1],
			})
		}
		if IsValidLinkWithKeyword(row[10], "") {
			socialMedias = append(socialMedias, hoteldto.SocialMedia{
				Platform: "Website",
				Link:     row[10],
			})
		}
		if IsValidLinkWithKeyword(row[11], "instagram") {
			socialMedias = append(socialMedias, hoteldto.SocialMedia{
				Platform: "Instagram",
				Link:     row[11],
			})
		}

		jsonSocialMedias, err = json.Marshal(socialMedias)
		if err != nil {
			logger.Error(ctx, "failed to marshal social media: %w", err)
		}

		reqHotel := &hoteldto.CreateHotelRequest{
			Name:         row[0],
			SubDistrict:  row[1],
			District:     row[2],
			Email:        row[3],
			Province:     row[4],
			Description:  row[5],
			Rating:       rating,
			Facilities:   strings.Split(row[8], ","),
			NearbyPlaces: string(jsonNearbyPlaces),
			SocialMedias: string(jsonSocialMedias),
		}

		// 5️⃣ Transformasi ke entity Hotel
		if _, err := hu.CreateHotel(ctx, reqHotel); err != nil {
			return fmt.Errorf("failed to create hotel at row %d: %w", i, err)
		}
	}

	return nil

}

func IsValidLinkWithKeyword(s, keyword string) bool {
	if s == "" {
		return false
	}

	// cek valid URL
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}

	// jika keyword kosong, cukup return true
	if keyword == "" {
		return true
	}

	// cek apakah string mengandung keyword (case-insensitive)
	return strings.Contains(strings.ToLower(s), strings.ToLower(keyword))
}
