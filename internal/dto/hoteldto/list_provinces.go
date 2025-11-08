package hoteldto

import "wtm-backend/internal/dto"

type ListProvincesRequest struct {
	dto.PaginationRequest `json:",inline"`
}
type ListProvincesResponse struct {
	Provinces []string `json:"provinces"`
	Total     int64    `json:"total"`
}
