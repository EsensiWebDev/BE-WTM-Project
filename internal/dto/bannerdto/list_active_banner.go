package bannerdto

import (
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto"
)

type ListBannerRequest struct {
	dto.PaginationRequest `json:",inline"`
	IsActive              bool `json:"is_active"`
}

type ListBannerResponse struct {
	Banners []entity.Banner `json:"banners"`
	Total   int64           `json:"total"`
}
