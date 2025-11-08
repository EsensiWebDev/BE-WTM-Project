package entity

import "time"

type EmailLog struct {
	ID              uint             `json:"id"`
	To              string           `json:"to"`
	Subject         string           `json:"subject"`
	Body            string           `json:"body"`
	Meta            MetadataEmailLog `json:"meta"`
	EmailTemplateID uint             `json:"email_template_id"`
	StatusID        uint             `json:"status_id"`
	CreatedAt       time.Time        `json:"created_at"`
}

type MetadataEmailLog struct {
	HotelID   string `json:"hotel_id"`
	HotelName string `json:"hotel_name"`
	AgentID   string `json:"agent_id"`
	Notes     string `json:"notes"`
}
