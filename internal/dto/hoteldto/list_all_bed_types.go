package hoteldto

import "wtm-backend/internal/dto"

type ListAllBedTypesRequest struct {
	dto.PaginationRequest `json:",inline"`
}

type ListAllBedTypesResponse struct {
	BedTypes []string `json:"bed_types"`
	Total    int64    `json:"total"`
}
