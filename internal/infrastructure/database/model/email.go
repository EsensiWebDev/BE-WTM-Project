package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type EmailTemplate struct {
	gorm.Model
	ExternalID       ExternalID `gorm:"embedded"`
	Name             string     `gorm:"type:varchar(100);uniqueIndex;not null"`
	Subject          string     `gorm:"type:varchar(255);not null"`
	Body             string     `gorm:"type:text;not null"`
	IsSignatureImage bool       `gorm:"type:boolean;not null"`
	Signature        string     `gorm:"type:text;not null"`
}

func (b *EmailTemplate) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}

type EmailLog struct {
	gorm.Model
	ExternalID      ExternalID     `gorm:"embedded"`
	To              string         `gorm:"type:text"`
	Subject         string         `gorm:"type:text"`
	Body            string         `gorm:"type:text"`
	Meta            datatypes.JSON `gorm:"type:jsonb"`
	EmailTemplateID uint           `gorm:"not null"`
	StatusID        uint           `gorm:"not null"`

	EmailStatus   StatusEmail   `gorm:"foreignKey:StatusID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	EmailTemplate EmailTemplate `gorm:"foreignKey:EmailTemplateID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (b *EmailLog) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}

type StatusEmail struct {
	ID         uint       `gorm:"primaryKey"` // override default gorm.Model ID
	Status     string     `gorm:"type:varchar(50);uniqueIndex;not null"`
	ExternalID ExternalID `gorm:"embedded"`
}

func (b *StatusEmail) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}
