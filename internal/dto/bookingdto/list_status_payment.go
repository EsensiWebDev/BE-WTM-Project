package bookingdto

import "wtm-backend/internal/domain/entity"

type ListStatusPaymentResponse struct {
	Data []entity.StatusPayment `json:"payments"`
}
