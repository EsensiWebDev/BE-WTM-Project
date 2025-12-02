package bannerdto

type DetailBannerRequest struct {
	BannerID string `json:"id" uri:"id"`
}
type DetailBannerResponse struct {
	Banner BannerData `json:"banner"`
}
