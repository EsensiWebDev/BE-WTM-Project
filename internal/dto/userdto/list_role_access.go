package userdto

type ListRoleAccessResponse struct {
	Role   string                     `json:"role"`
	Access map[string]map[string]bool `json:"access"` // page → action → allowed
}
