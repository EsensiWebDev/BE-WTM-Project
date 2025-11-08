package userdto

import (
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto"
)

type ListStatusUsersRequest struct {
	dto.PaginationRequest `json:",inline"`
}

type ListStatusUsersResponse struct {
	StatusUsers []entity.StatusUser `json:"status_users"`
}
