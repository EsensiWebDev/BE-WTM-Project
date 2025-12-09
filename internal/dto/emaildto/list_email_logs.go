package emaildto

import "wtm-backend/internal/dto"

type ListEmailLogsRequest struct {
	dto.PaginationRequest `json:",inline"`
}

type ListEmailLogsResponse struct {
	EmailLogs []EmailLogResponse `json:"email_logs"`
	Total     int64              `json:"total"`
}

// EmailLogResponse represents a single email log entry in the response
type EmailLogResponse struct {
	DateTime  string `json:"date_time"`
	EmailType string `json:"email_type"`
	HotelName string `json:"hotel_name"`
	Status    string `json:"status"`
	Notes     string `json:"notes"`
}
