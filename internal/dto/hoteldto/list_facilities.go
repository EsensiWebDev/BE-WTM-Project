package hoteldto

import (
	"wtm-backend/internal/dto"
)

type ListFacilitiesRequest struct {
	dto.PaginationRequest `json:",inline"`
}

type ListFacilitiesResponse struct {
	Facilities []string `json:"facilities"`
	Total      int64    `json:"total"`
}
