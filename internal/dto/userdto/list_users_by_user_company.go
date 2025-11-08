package userdto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"wtm-backend/internal/dto"
)

type ListUsersByAgentCompanyRequest struct {
	ID                    uint `json:"id" form:"id"`
	dto.PaginationRequest `json:",inline"`
}

type ListUsersByAgentCompanyResponse struct {
	Users []ListUsersByAgentCompanyData `json:"users"`
}

type ListUsersByAgentCompanyData struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (r *ListUsersByAgentCompanyRequest) Validate() error {
	return validation.ValidateStruct(r, validation.Field(&r.ID, validation.Required))
}
