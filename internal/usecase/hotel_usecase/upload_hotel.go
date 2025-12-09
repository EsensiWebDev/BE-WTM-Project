package hotel_usecase

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/mail"
	"sort"
	"strconv"
	"strings"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/pkg/logger"
)

func (hu *HotelUsecase) UploadHotel(ctx context.Context, req *hoteldto.UploadHotelRequest) (bool, error) {
	// 1. Buka file CSV
	file, err := req.File.Open()
	if err != nil {
		return false, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Baca semua konten untuk handle Excel artifacts
	content, err := io.ReadAll(file)
	if err != nil {
		return false, fmt.Errorf("failed to read file: %w", err)
	}

	// Clean CSV dari Excel artifacts
	cleanedContent := string(content)
	cleanedContent = strings.ReplaceAll(cleanedContent, ";;;;;;;;;;", "") // Hapus semicolon berlebih
	cleanedContent = strings.TrimSpace(cleanedContent)

	reader := csv.NewReader(strings.NewReader(cleanedContent))
	reader.Comma = ';'
	reader.Comment = '#'
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		return false, fmt.Errorf("failed to read csv: %w", err)
	}

	// Cari header row (skip baris yang hanya komentar atau kosong)
	headerIndex := -1
	for i, row := range records {
		if len(row) > 0 && !strings.HasPrefix(strings.TrimSpace(row[0]), "#") {
			// Cek apakah ini header (contains column names)
			rowStr := strings.ToLower(strings.Join(row, " "))
			if strings.Contains(rowStr, "name") && strings.Contains(rowStr, "email") {
				headerIndex = i
				break
			}
		}
	}

	if headerIndex == -1 {
		return false, fmt.Errorf("header row not found")
	}

	headers := records[headerIndex]

	// Validasi header minimal
	requiredHeaders := []string{"name", "email"}
	headerMap := make(map[string]int)
	for i, h := range headers {
		key := strings.ToLower(strings.TrimSpace(h))
		headerMap[key] = i
	}

	for _, req := range requiredHeaders {
		if _, exists := headerMap[req]; !exists {
			return false, fmt.Errorf("missing required header: %s", req)
		}
	}

	// Track hasil
	successCount := 0
	failedRows := make(map[int][]string) // row number -> errors

	// Process data rows
	for i := headerIndex + 1; i < len(records); i++ {
		row := records[i]
		rowNum := i + 1 // 1-based

		if len(row) == 0 || strings.HasPrefix(strings.TrimSpace(row[0]), "#") {
			continue
		}

		// Parse dan validasi
		hotel, errs := hu.parseAndValidateHotelRow(row, rowNum, headerMap)
		if len(errs) > 0 {
			failedRows[rowNum] = errs
			continue
		}

		// Create hotel
		if _, err := hu.CreateHotel(ctx, hotel); err != nil {
			failedRows[rowNum] = []string{fmt.Sprintf("failed to create hotel: %v", err)}
			continue
		}

		successCount++
	}

	// Generate hasil
	if len(failedRows) > 0 {
		// Buat error message yang ringkas
		var failedRowNumbers []int
		for rowNum := range failedRows {
			failedRowNumbers = append(failedRowNumbers, rowNum)
		}
		sort.Ints(failedRowNumbers)

		// Log detail error
		logger.Warn(ctx, "Upload completed with %d successes, %d failures",
			successCount, len(failedRows))

		for _, rowNum := range failedRowNumbers {
			errors := failedRows[rowNum]
			logger.Warn(ctx, "Row %d errors: %v", rowNum, errors)
		}

		// Return error yang ringkas
		if successCount == 0 {
			return false, fmt.Errorf("all rows failed. Failed rows: %v", failedRowNumbers)
		}

		return true, fmt.Errorf("%d rows succeeded, %d rows failed. Failed rows: %v",
			successCount, len(failedRows), failedRowNumbers)
	}

	logger.Info(ctx, "Upload completed successfully: %d hotels created", successCount)
	return true, nil
}

func (hu *HotelUsecase) parseAndValidateHotelRow(row []string, rowNum int, headerMap map[string]int) (*hoteldto.CreateHotelRequest, []string) {
	var errors []string

	// Helper untuk get value dari header map
	getValue := func(key string) string {
		if idx, exists := headerMap[key]; exists && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}

		// Coba dengan alternative names
		if key == "sub_district" {
			if idx, exists := headerMap["subdistrict"]; exists && idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
		} else if key == "nearby_places" {
			if idx, exists := headerMap["nearbyplaces"]; exists && idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
		}

		return ""
	}

	// Validasi required fields
	name := getValue("name")
	if name == "" {
		errors = append(errors, "name is required")
	}

	email := getValue("email")
	if email == "" {
		errors = append(errors, "email is required")
	} else if !isValidEmail(email) {
		errors = append(errors, fmt.Sprintf("invalid email: %s", email))
	}

	// Parse rating
	rating := 0
	ratingStr := getValue("rating")
	if ratingStr != "" {
		if val, err := strconv.Atoi(ratingStr); err == nil {
			if val >= 0 && val <= 5 {
				rating = val
			} else {
				errors = append(errors, fmt.Sprintf("rating must be 0-5, got %d", val))
			}
		} else {
			errors = append(errors, fmt.Sprintf("invalid rating: %s", ratingStr))
		}
	}

	// Parse nearby places
	nearbyPlaces, nearbyErrs := parseNearbyPlacesWithErrors(getValue("nearby_places"))
	if len(nearbyErrs) > 0 {
		errors = append(errors, nearbyErrs...)
	}

	jsonNearbyPlaces, err := json.Marshal(nearbyPlaces)
	if err != nil {
		errors = append(errors, "failed to marshal nearby places")
		jsonNearbyPlaces = []byte("[]")
	}

	// Parse facilities
	var facilities []string
	if facilitiesStr := getValue("facilities"); facilitiesStr != "" {
		facilities = parseFacilities(facilitiesStr)
	}

	// Parse social medias
	socialMedias := parseSocialMedias(
		getValue("tiktok"),
		getValue("website"),
		getValue("instagram"),
	)

	jsonSocialMedias, err := json.Marshal(socialMedias)
	if err != nil {
		errors = append(errors, "failed to marshal social medias")
		jsonSocialMedias = []byte("[]")
	}

	// Jika ada critical errors, return nil
	if len(errors) > 0 && (name == "" || email == "" || !isValidEmail(email)) {
		return nil, errors
	}

	// Build hotel request (masih return meski ada non-critical errors)
	return &hoteldto.CreateHotelRequest{
		Name:         name,
		SubDistrict:  getValue("sub_district"),
		District:     getValue("district"),
		Email:        email,
		Province:     getValue("province"),
		Description:  getValue("description"),
		Rating:       rating,
		Facilities:   facilities,
		NearbyPlaces: string(jsonNearbyPlaces),
		SocialMedias: string(jsonSocialMedias),
	}, errors
}

// Update parseNearbyPlaces untuk return errors
func parseNearbyPlacesWithErrors(nearbyStr string) ([]hoteldto.NearbyPlace, []string) {
	if strings.TrimSpace(nearbyStr) == "" {
		return []hoteldto.NearbyPlace{}, nil
	}

	var places []hoteldto.NearbyPlace
	var errors []string
	placeEntries := strings.Split(nearbyStr, "|")

	for i, entry := range placeEntries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		parts := strings.Split(entry, ",")
		if len(parts) != 2 {
			errors = append(errors, fmt.Sprintf("nearby place entry %d: invalid format '%s'", i+1, entry))
			continue
		}

		name := strings.TrimSpace(parts[0])
		distanceStr := strings.TrimSpace(parts[1])

		distance, err := strconv.ParseFloat(distanceStr, 64)
		if err != nil {
			errors = append(errors, fmt.Sprintf("nearby place entry %d: invalid distance '%s'", i+1, distanceStr))
			continue
		}

		places = append(places, hoteldto.NearbyPlace{
			Name:     name,
			Distance: distance,
		})
	}

	return places, errors
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
