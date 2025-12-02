package model

import (
	"gorm.io/gorm"
	"time"
)

type StatusUser struct {
	ID         uint       `gorm:"primaryKey"` // override default gorm.Model ID
	Status     string     `json:"status"`
	ExternalID ExternalID `gorm:"embedded"`
}

type AgentCompany struct {
	gorm.Model
	Name       string     `json:"name"`
	ExternalID ExternalID `gorm:"embedded"`
}

func (b *AgentCompany) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}

type User struct {
	gorm.Model
	FullName       string     `json:"full_name"`
	AgentCompanyID *uint      `json:"agent_company_id" gorm:"index"`
	Username       string     `json:"username" gorm:"uniqueIndex:idx_users_username_active,where:deleted_at IS NULL;not null"`
	Password       string     `json:"password"`
	StatusID       uint       `json:"status_id" gorm:"index; default:1"`
	RoleID         uint       `json:"role_id" gorm:"index"`
	PromoGroupID   *uint      `json:"promo_group_id" gorm:"index"`
	Email          string     `json:"email" gorm:"uniqueIndex:idx_users_email_active,where:deleted_at IS NULL;not null"`
	Phone          string     `json:"phone" gorm:"uniqueIndex:idx_users_phone_active,where:deleted_at IS NULL;not null"`
	KakaoTalkID    string     `json:"kakao_talk_id"`
	Certificate    string     `json:"certificate"`
	NameCard       string     `json:"name_card"`
	PhotoSelfie    string     `json:"photo_selfie"`
	PhotoIDCard    string     `json:"photo_id_card"`
	ExternalID     ExternalID `gorm:"embedded"`

	Status       StatusUser    `gorm:"foreignKey:StatusID"`
	AgentCompany *AgentCompany `gorm:"foreignKey:AgentCompanyID"`

	Role *Role `gorm:"foreignKey:RoleID"`

	PromoGroup *PromoGroup `gorm:"foreignKey:PromoGroupID"`

	UserNotificationSettings []UserNotificationSetting `gorm:"foreignKey:UserID"`
}

func (b *User) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}

type Role struct {
	gorm.Model
	ExternalID  ExternalID   `gorm:"embedded"`
	Role        string       `json:"role"`
	Permissions []Permission `gorm:"many2many:role_permissions"` // many-to-many
}

type RolePermission struct {
	gorm.Model
	ExternalID   ExternalID `gorm:"embedded"`
	RoleID       uint       `json:"role_id" gorm:"index"`
	PermissionID uint       `json:"permission_id" gorm:"index"`

	Role       Role       `gorm:"foreignKey:RoleID"`
	Permission Permission `gorm:"foreignKey:PermissionID"`
}

func (b *RolePermission) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}

type Permission struct {
	gorm.Model
	ExternalID ExternalID `gorm:"embedded"`
	Permission string     `json:"permission"` // e.g., "account:create"
	Page       string     `json:"page"`       // optional
	Action     string     `json:"action"`     // optional

	Role []Role `gorm:"many2many:role_permissions"` // many-to-many
}

type PasswordResetToken struct {
	gorm.Model
	ExternalID ExternalID `gorm:"embedded"`
	UserID     uint       `json:"user_id" gorm:"index"`
	Token      string     `json:"token" gorm:"uniqueIndex;not null"`
	ExpiresAt  time.Time  `json:"expires_at"`
	Used       bool       `json:"used" gorm:"default:false"`
}

func (b *PasswordResetToken) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}
