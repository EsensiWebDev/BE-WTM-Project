package bookingdto

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type AddToCartRequest struct {
	RoomPriceID           uint   `json:"room_price_id"`
	CheckInDate           string `json:"check_in_date"`
	CheckOutDate          string `json:"check_out_date"`
	Quantity              int    `json:"quantity"`
	RoomTypeAdditionalIDs []uint `json:"room_type_additional_ids"`
	PromoID               uint   `json:"promo_id"`
}

func (r *AddToCartRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.RoomPriceID, validation.Required.Error("Room Price Id is required")),
		validation.Field(&r.CheckInDate, validation.Required.Error("Check In Date is required")),
		validation.Field(&r.CheckOutDate, validation.Required.Error("Check Out Date is required")),
		validation.Field(&r.Quantity, validation.Required.Error("Quantity is required")),
	)
}
