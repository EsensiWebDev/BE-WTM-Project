package hotel_usecase

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/mail"
	"strconv"
	"strings"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/pkg/logger"
)

func (hu *HotelUsecase) UploadHotel(ctx context.Context, req *hoteldto.UploadHotelRequest) error {
	// 1. Buka file CSV
	file, err := req.File.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.Comment = '#'        // Otomatis skip baris yang diawali #
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read csv: %w", err)
	}

	// 2. Validasi minimal baris
	if len(records) < 2 {
		return fmt.Errorf("csv file must have at least header and one data row")
	}

	// 3. Map header untuk validasi
	headers := records[0]
	expectedHeaders := []string{
		"name", "sub_district", "district", "email", "province",
		"description", "rating", "nearby_places", "facilities",
		"tiktok", "website", "instagram",
	}

	if len(headers) != len(expectedHeaders) {
		return fmt.Errorf("invalid header count. expected %d, got %d",
			len(expectedHeaders), len(headers))
	}

	// 4. Iterasi data rows (skip header)
	successCount := 0
	failedRows := []string{}

	for i, row := range records[1:] {
		rowNum := i + 2 // +1 untuk header, +1 untuk 1-based index

		// Skip baris kosong atau komentar (dengan Comment sudah handled)
		if len(row) == 0 {
			continue
		}

		// Validasi row length
		if len(row) < len(expectedHeaders) {
			failedRows = append(failedRows, fmt.Sprintf("Row %d: insufficient columns", rowNum))
			continue
		}

		// 5. Parse data dengan error handling per field
		hotel, err := hu.parseHotelRow(ctx, row, rowNum)
		if err != nil {
			failedRows = append(failedRows, fmt.Sprintf("Row %d: %v", rowNum, err))
			continue
		}

		// 6. Create hotel
		if _, err := hu.CreateHotel(ctx, hotel); err != nil {
			failedRows = append(failedRows, fmt.Sprintf("Row %d: failed to create hotel - %v", rowNum, err))
			continue
		}

		successCount++
	}

	// 7. Summary
	if len(failedRows) > 0 {
		logger.Warn(ctx, "Upload completed with %d successes, %d failures",
			successCount, len(failedRows))
		for _, msg := range failedRows {
			logger.Warn(ctx, msg)
		}

		if successCount == 0 {
			return fmt.Errorf("all rows failed to process")
		}
		return fmt.Errorf("processed with %d errors", len(failedRows))
	}

	logger.Info(ctx, "Upload completed successfully: %d hotels created", successCount)
	return nil
}

// Helper function untuk parse row
func (hu *HotelUsecase) parseHotelRow(ctx context.Context, row []string, rowNum int) (*hoteldto.CreateHotelRequest, error) {
	// Map column indexes
	const (
		colName = iota
		colSubDistrict
		colDistrict
		colEmail
		colProvince
		colDescription
		colRating
		colNearbyPlaces
		colFacilities
		colTiktok
		colWebsite
		colInstagram
	)

	// 1. Parse rating
	var rating int
	if strings.TrimSpace(row[colRating]) != "" {
		if _, err := fmt.Sscanf(row[colRating], "%d", &rating); err != nil {
			logger.Warn(ctx, "Row %d: invalid rating '%s', using default 0",
				rowNum, row[colRating])
			rating = 0
		}
	}

	// Validate rating range
	if rating < 0 || rating > 5 {
		rating = 5
	}

	// 2. Parse nearby places
	nearbyPlaces, err := parseNearbyPlaces(row[colNearbyPlaces])
	if err != nil {
		logger.Warn(ctx, "Row %d: failed to parse nearby places: %v", rowNum, err)
	}

	jsonNearbyPlaces, err := json.Marshal(nearbyPlaces)
	if err != nil {
		logger.Warn(ctx, "Row %d: failed to marshal nearby places: %v", rowNum, err)
		jsonNearbyPlaces = []byte("[]")
	}

	// 3. Parse facilities
	var facilities []string
	if strings.TrimSpace(row[colFacilities]) != "" {
		facilities = parseFacilities(row[colFacilities])
	}

	// 4. Parse social medias
	socialMedias := parseSocialMedias(
		row[colTiktok],
		row[colWebsite],
		row[colInstagram],
	)

	jsonSocialMedias, err := json.Marshal(socialMedias)
	if err != nil {
		logger.Warn(ctx, "Row %d: failed to marshal social medias: %v", rowNum, err)
		jsonSocialMedias = []byte("[]")
	}

	// 5. Validasi required fields
	if strings.TrimSpace(row[colName]) == "" {
		return nil, fmt.Errorf("name is required")
	}
	if strings.TrimSpace(row[colEmail]) == "" {
		return nil, fmt.Errorf("email is required")
	}
	email := strings.TrimSpace(row[colEmail])
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	if !isValidEmail(email) {
		return nil, fmt.Errorf("invalid email format: %s", email)
	}

	// 6. Build request
	return &hoteldto.CreateHotelRequest{
		Name:         strings.TrimSpace(row[colName]),
		SubDistrict:  strings.TrimSpace(row[colSubDistrict]),
		District:     strings.TrimSpace(row[colDistrict]),
		Email:        email,
		Province:     strings.TrimSpace(row[colProvince]),
		Description:  strings.TrimSpace(row[colDescription]),
		Rating:       rating,
		Facilities:   facilities,
		NearbyPlaces: string(jsonNearbyPlaces),
		SocialMedias: string(jsonSocialMedias),
	}, nil
}

func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}

	// Cek format dengan net/mail
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Helper untuk parse nearby places
func parseNearbyPlaces(nearbyStr string) ([]hoteldto.NearbyPlace, error) {
	if strings.TrimSpace(nearbyStr) == "" {
		return []hoteldto.NearbyPlace{}, nil
	}

	var places []hoteldto.NearbyPlace
	placeEntries := strings.Split(nearbyStr, "|")

	for _, entry := range placeEntries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		parts := strings.Split(entry, ",")
		if len(parts) != 2 {
			continue // Skip invalid format
		}

		name := strings.TrimSpace(parts[0])
		distanceStr := strings.TrimSpace(parts[1])

		distance, err := strconv.ParseFloat(distanceStr, 64)
		if err != nil {
			continue // Skip invalid distance
		}

		places = append(places, hoteldto.NearbyPlace{
			Name:     name,
			Distance: distance,
		})
	}

	return places, nil
}

// Helper untuk parse facilities
func parseFacilities(facilitiesStr string) []string {
	if strings.TrimSpace(facilitiesStr) == "" {
		return []string{}
	}

	var facilities []string
	items := strings.Split(facilitiesStr, ",")

	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			facilities = append(facilities, item)
		}
	}

	return facilities
}

// Helper untuk parse social medias
func parseSocialMedias(tiktok, website, instagram string) []hoteldto.SocialMedia {
	var socialMedias []hoteldto.SocialMedia

	// Tiktok
	if IsValidLinkWithKeyword(tiktok, "tiktok") {
		socialMedias = append(socialMedias, hoteldto.SocialMedia{
			Platform: "Tiktok",
			Link:     strings.TrimSpace(tiktok),
		})
	}

	// Website
	if IsValidLinkWithKeyword(website, "") {
		socialMedias = append(socialMedias, hoteldto.SocialMedia{
			Platform: "Website",
			Link:     strings.TrimSpace(website),
		})
	}

	// Instagram
	if IsValidLinkWithKeyword(instagram, "instagram") {
		socialMedias = append(socialMedias, hoteldto.SocialMedia{
			Platform: "Instagram",
			Link:     strings.TrimSpace(instagram),
		})
	}

	return socialMedias
}

// Helper untuk validasi link
func IsValidLinkWithKeyword(link, keyword string) bool {
	link = strings.TrimSpace(link)
	if link == "" || !strings.HasPrefix(link, "http") {
		return false
	}

	if keyword == "" {
		return true
	}

	return strings.Contains(strings.ToLower(link), strings.ToLower(keyword))
}
