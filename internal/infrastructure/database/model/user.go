package model

import (
	"gorm.io/gorm"
	"time"
)

type StatusUser struct {
	ID     uint   `gorm:"primaryKey"` // override default gorm.Model ID
	Status string `json:"status"`
}

type AgentCompany struct {
	gorm.Model
	Name string `json:"name"`
}

type User struct {
	gorm.Model
	FullName       string `json:"full_name"`
	AgentCompanyID *uint  `json:"agent_company_id" gorm:"index"`
	Username       string `json:"username" gorm:"uniqueIndex;not null"`
	Password       string `json:"password"`
	StatusID       uint   `json:"status_id" gorm:"index; default:1"`
	RoleID         uint   `json:"role_id" gorm:"index"`
	PromoGroupID   *uint  `json:"promo_group_id" gorm:"index"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	KakaoTalkID    string `json:"kakao_talk_id"`
	Certificate    string `json:"certificate"`
	NameCard       string `json:"name_card"`
	PhotoSelfie    string `json:"photo_selfie"`
	PhotoIDCard    string `json:"photo_id_card"`

	Status       StatusUser    `gorm:"foreignKey:StatusID"`
	AgentCompany *AgentCompany `gorm:"foreignKey:AgentCompanyID"`

	Role *Role `gorm:"foreignKey:RoleID"`

	PromoGroup *PromoGroup `gorm:"foreignKey:PromoGroupID"`

	UserNotificationSettings []UserNotificationSetting `gorm:"foreignKey:UserID"`
}

type Role struct {
	gorm.Model
	Role        string       `json:"role"`
	Permissions []Permission `gorm:"many2many:role_permissions"` // many-to-many
}

type RolePermission struct {
	gorm.Model
	RoleID       uint `json:"role_id" gorm:"index"`
	PermissionID uint `json:"permission_id" gorm:"index"`

	Role       Role       `gorm:"foreignKey:RoleID"`
	Permission Permission `gorm:"foreignKey:PermissionID"`
}

type Permission struct {
	gorm.Model
	Permission string `json:"permission"` // e.g., "account:create"
	Page       string `json:"page"`       // optional
	Action     string `json:"action"`     // optional

	Role []Role `gorm:"many2many:role_permissions"` // many-to-many
}

type PasswordResetToken struct {
	gorm.Model
	UserID    uint      `json:"user_id" gorm:"index"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used" gorm:"default:false"`
}
