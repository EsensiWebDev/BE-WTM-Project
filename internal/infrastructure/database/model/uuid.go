package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExternalID struct {
	ExternalID string `json:"external_id" gorm:"column:external_id;not null;uniqueIndex"`
}

func (e *ExternalID) BeforeCreate(tx *gorm.DB) error {
	if e.ExternalID == "" {
		uuidString, err := uuid.NewV7()
		if err != nil {
			uuidString = uuid.New() // fallback ke v4
		}
		e.ExternalID = uuidString.String()
	}
	return nil
}
