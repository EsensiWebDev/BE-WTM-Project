package userdto

import validation "github.com/go-ozzo/ozzo-validation"

type UpdateRoleAccessRequest struct {
	Role    string `json:"role"`
	Page    string `json:"page"`
	Action  string `json:"action"`
	Allowed bool   `json:"allowed"`
}

func (r *UpdateRoleAccessRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Role, validation.Required, validation.In("admin", "support", "agent").Error("Role must be one of: admin, support, agent")),
		validation.Field(&r.Page, validation.Required.Error("Page is required")),
		validation.Field(&r.Action, validation.Required.Error("Action is required")),
	)
}
