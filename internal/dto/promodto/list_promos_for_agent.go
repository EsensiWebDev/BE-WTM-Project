package promodto

import "wtm-backend/internal/dto"

type ListPromosForAgentRequest struct {
	dto.PaginationRequest `json:",inline"`
}

type ListPromosForAgentResponse struct {
	Data  []PromosForAgent `json:"data"`
	Total int64            `json:"total"`
}

type PromosForAgent struct {
	ID          uint     `json:"id"`
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Hotel       []string `json:"hotel"`
}
