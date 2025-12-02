package model

import (
	"gorm.io/gorm"
	"time"
)

type Notification struct {
	gorm.Model
	ExternalID  ExternalID `gorm:"embedded"`
	UserID      uint       `gorm:"index;not null"`
	Title       string     `gorm:"type:varchar(255);not null"`
	Message     string     `gorm:"type:text;not null"`
	RedirectURL string     `gorm:"type:varchar(255)"`
	IsRead      bool       `gorm:"type:boolean;not null"`
	ReadAt      time.Time
}

func (b *Notification) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}

type UserNotificationSetting struct {
	gorm.Model
	ExternalID ExternalID `gorm:"embedded"`
	UserID     uint       `gorm:"index;not null"`
	Channel    string     `gorm:"type:varchar(20);not null"` // "email", "web"
	Type       string     `gorm:"type:varchar(50);not null"` // "booking", "reject", "both"
	IsEnabled  bool       `gorm:"not null"`

	User User `gorm:"foreignKey:UserID"`
}

func (b *UserNotificationSetting) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}
