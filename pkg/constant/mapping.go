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
	StatusBookingConfirmedID:       StatusBookingConfirmed,
	StatusBookingWaitingApprovalID: StatusBookingWaitingApproval,
	StatusBookingRejectedID:        StatusBookingRejected,
}

// slice untuk urutan
var StatusBookingOrder = []int{
	StatusBookingConfirmedID,
	StatusBookingWaitingApprovalID,
	StatusBookingRejectedID,
}

var MapStatusPayment = map[int]string{
	StatusPaymentUnpaidID: StatusPaymentUnpaid,
	StatusPaymentPaidID:   StatusPaymentPaid,
}

// slice untuk urutan
var StatusPaymentOrder = []int{
	StatusPaymentUnpaidID,
	StatusPaymentPaidID,
}
