package bookingdto

type CheckOutCartRequest struct {
	BookingID uint             `json:"booking_id"`
	Guests    []string         `json:"guests"`
	Details   []CheckOutDetail `json:"detail"`
}

type CheckOutDetail struct {
	BookingDetailID uint   `json:"booking_detail_id"`
	Guest           string `json:"guest"`
}
