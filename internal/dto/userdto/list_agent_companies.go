package userdto

import (
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto"
)

type ListAgentCompaniesRequest struct {
	dto.PaginationRequest `json:",inline"`
}

type ListAgentCompaniesResponse struct {
	AgentCompanies []entity.AgentCompany `json:"agent_companies"`
}
