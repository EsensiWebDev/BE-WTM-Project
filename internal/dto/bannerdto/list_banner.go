package bannerdto

import "wtm-backend/internal/dto"

type ListBannerRequest struct {
	dto.PaginationRequest `json:",inline"`
	IsActive              *bool `json:"is_active"`
}

type ListBannerResponse struct {
	Banners []BannerData `json:"banners"`
	Total   int64        `json:"total"`
}

type BannerData struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}
