package bannerdto

type ListActiveBannerResponse struct {
	Banners []ActiveBanner
}

// ActiveBanner is a struct that represents an active banner
type ActiveBanner struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
}
