package userdto

// ProfileResponse represents the structure of the user profile response.
type ProfileResponse struct {
	ID                  uint                  `json:"id"`
	Username            string                `json:"username"`
	Password            string                `json:"password"`
	FullName            string                `json:"full_name"`
	Email               string                `json:"email"`
	Phone               string                `json:"phone"`
	PhotoProfile        string                `json:"photo_profile"`
	Certificate         string                `json:"certificate,omitempty"`
	NameCard            string                `json:"name_card,omitempty"`
	KakaoTalkID         string                `json:"kakao_talk_id,omitempty"`
	Status              string                `json:"status,omitempty"`
	AgentCompany        string                `json:"agent_company,omitempty"`
	Currency            string                `json:"currency,omitempty"`
	NotificationSetting []NotificationSetting `json:"notification_settings,omitempty"`
}

type NotificationSetting struct {
	Channel  string `json:"channel"`
	Type     string `json:"type"`
	IsEnable bool   `json:"is_enable"`
}
