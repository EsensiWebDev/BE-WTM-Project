package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type EmailTemplate struct {
	gorm.Model
	Name             string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Subject          string `gorm:"type:varchar(255);not null"`
	Body             string `gorm:"type:text;not null"`
	IsSignatureImage bool   `gorm:"type:boolean;not null"`
	Signature        string `gorm:"type:text;not null"`
}

type EmailLog struct {
	gorm.Model
	To              string         `gorm:"type:varchar(255);not null"`
	Subject         string         `gorm:"type:varchar(255);not null"`
	Body            string         `gorm:"type:text;not null"`
	Meta            datatypes.JSON `gorm:"type:jsonb;not null"`
	EmailTemplateID uint           `gorm:"not null"`
	StatusID        uint           `gorm:"not null"`

	EmailStatus StatusEmail `gorm:"foreignKey:StatusID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type StatusEmail struct {
	ID     uint   `gorm:"primaryKey"` // override default gorm.Model ID
	Status string `gorm:"type:varchar(50);uniqueIndex;not null"`
}
