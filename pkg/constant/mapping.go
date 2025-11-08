package constant

var MapStatusHotel = map[int]string{
	StatusHotelApprovedID: StatusHotelApproved,
	StatusHotelInReviewID: StatusHotelInReview,
	StatusHotelRejectedID: StatusHotelRejected,
}

var MapStatusEmailLog = map[int]string{
	StatusEmailPendingID: StatusEmailPending,
	StatusEmailSuccessID: StatusEmailSuccess,
	StatusEmailFailedID:  StatusEmailFailed,
}

var MapStatusBooking = map[int]string{
	StatusBookingApprovedID: StatusBookingApproved,
	StatusBookingInReviewID: StatusBookingInReview,
	StatusBookingRejectedID: StatusBookingRejected,
}

var MapStatusPayment = map[int]string{
	StatusPaymentUnpaidID: StatusPaymentUnpaid,
	StatusPaymentPaidID:   StatusPaymentPaid,
}
