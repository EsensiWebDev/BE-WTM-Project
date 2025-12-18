package entity

type User struct {
	FullName       string
	AgentCompanyID *uint
	Username       string
	Password       string
	StatusID       uint
	RoleID         uint
	PromoGroupID   *uint
	Email          string
	Phone          string
	KakaoTalkID    string
	Certificate    string
	NameCard       string
	PhotoSelfie    string
	PhotoIDCard    string
	Currency       string // Agent currency preference (set by admin)

	//additional fields
	ID                       uint
	ExternalID               string
	StatusName               string
	RoleName                 string
	Permissions              []string
	AgentCompanyName         string
	PromoGroupName           string
	UserNotificationSettings []UserNotificationSetting
}

type UserNotificationSetting struct {
	UserID    uint
	Channel   string
	Type      string
	IsEnabled bool
}

type AgentCompany struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// UserMin is a minimal representation of a user, typically used for listings or summaries.
type UserMin struct {
	ID          uint     `json:"id"`
	Username    string   `json:"username"`
	RoleID      uint     `json:"role_id"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	PhotoURL    string   `json:"photo_url"`
	FullName    string   `json:"full_name"`
}

type Role struct {
	ID          uint
	Role        string
	Permissions []Permission
}

type Permission struct {
	ID         uint
	Permission string
	Page       string
	Action     string
}

type StatusUser struct {
	ID     uint   `json:"id"`
	Status string `json:"status"`
}
