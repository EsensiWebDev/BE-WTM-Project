package hoteldto

import "wtm-backend/internal/dto"

type ListAdditionalRoomsRequest struct {
	dto.PaginationRequest `json:",inline"`
}

type ListAdditionalRoomsResponse struct {
	AdditionalRooms []string `json:"additional_rooms"`
	Total           int64    `json:"total"`
}
