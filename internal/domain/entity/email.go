package entity

import "time"

type EmailLog struct {
	ID              uint              `json:"id"`
	To              string            `json:"to"`
	Subject         string            `json:"subject"`
	Body            string            `json:"body"`
	Meta            *MetadataEmailLog `json:"meta"`
	EmailTemplateID uint              `json:"email_template_id"`
	StatusID        uint              `json:"status_id"`
	CreatedAt       time.Time         `json:"created_at"`
	EmailType       string            `json:"email_type"`
}

type MetadataEmailLog struct {
	HotelName   string `json:"hotel_name"`
	AgentName   string `json:"agent_name"`
	Notes       string `json:"notes"`
	BookingCode string `json:"booking_code,omitempty"`
}
