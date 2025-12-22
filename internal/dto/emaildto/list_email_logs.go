package emaildto

import (
	"wtm-backend/internal/dto"
)

type ListEmailLogsRequest struct {
	dto.PaginationRequest `json:",inline"`
	Status                []string `form:"status" json:"status"`
	HotelName             string   `form:"hotel_name" json:"hotel_name"`
	BookingCode           string   `form:"booking_code" json:"booking_code"`
	DateFrom              string   `form:"date_from" json:"date_from"`
	DateTo                string   `form:"date_to" json:"date_to"`
}

type ListEmailLogsResponse struct {
	EmailLogs []EmailLogResponse `json:"email_logs"`
	Total     int64              `json:"total"`
}

// EmailLogResponse represents a single email log entry in the response
type EmailLogResponse struct {
	ID          uint   `json:"id"`
	DateTime    string `json:"date_time"`
	EmailType   string `json:"email_type"`
	HotelName   string `json:"hotel_name"`
	Status      string `json:"status"`
	Notes       string `json:"notes"`
	BookingCode string `json:"booking_code,omitempty"`
}

// EmailLogDetailResponse represents detailed email log information
type EmailLogDetailResponse struct {
	ID        uint   `json:"id"`
	To        string `json:"to"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
	DateTime  string `json:"date_time"`
	EmailType string `json:"email_type"`
	HotelName string `json:"hotel_name"`
	Status    string `json:"status"`
	Notes     string `json:"notes"`
}

type RetryEmailRequest struct {
	ID uint `json:"id" binding:"required"`
}

type RetryEmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
