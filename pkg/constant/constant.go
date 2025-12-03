package constant

const (
	ConstUser      = "user"
	ConstHotel     = "hotel"
	ConstBanner    = "banner"
	ConstPublic    = "public"
	ConstPrivate   = "private"
	ConstEmail     = "email"
	ConstSignature = "signature"
)

const (
	ConstBooking    = "booking"
	ConstSubBooking = "sub_booking"
	ConstPayment    = "payment"
	ConstReject     = "reject"
	ConstAll        = "all"
)

const (
	DefaultStatusSign     = 1
	DefaultRoleAgent      = 3
	DefaultRoleSuperAdmin = 1
)

const (
	ScopeControl    = "control"
	ScopeManagement = "management"
)

const (
	RoleAdmin         = "admin"
	RoleSupport       = "support"
	RoleAgent         = "agent"
	RoleSuperAdmin    = "super_admin"
	RoleSuperAdminID  = 1
	RoleAdminID       = 2
	RoleAgentID       = 3
	RoleSupportID     = 4
	RoleAdminCap      = "Admin"
	RoleSupportCap    = "Support"
	RoleSuperAdminCap = "Super Admin"
	RoleAgentCap      = "Agent"
)

const (
	StatusBookingInCart            = "In Cart"
	StatusBookingWaitingApproval   = "Waiting Approval"
	StatusBookingConfirmed         = "Confirmed"
	StatusBookingRejected          = "Rejected"
	StatusBookingCanceled          = "Canceled"
	StatusBookingInCartID          = 1
	StatusBookingWaitingApprovalID = 2
	StatusBookingConfirmedID       = 3
	StatusBookingRejectedID        = 4
	StatusBookingCanceledID        = 5
)

const (
	StatusPaymentUnpaid   = "Unpaid"
	StatusPaymentPaid     = "Paid"
	StatusPaymentUnpaidID = 1
	StatusPaymentPaidID   = 2
)

const (
	StatusHotelInReview   = "In Review"
	StatusHotelApproved   = "Approved"
	StatusHotelRejected   = "Rejected"
	StatusHotelInReviewID = 1
	StatusHotelApprovedID = 2
	StatusHotelRejectedID = 3
)

const (
	PromoTypeDiscount      = "Discount"
	PromoTypeFixedPrice    = "Fixed Price"
	PromoTypeRoomUpgrade   = "Room Upgrade"
	PromoTypeBenefit       = "Benefit"
	PromoTypeDiscountID    = 1
	PromoTypeFixedPriceID  = 2
	PromoTypeRoomUpgradeID = 3
	PromoTypeBenefitID     = 4
)

const (
	EmailAgentApproved       = "agent_approval"
	EmailAgentRejected       = "agent_rejection"
	EmailBookingConfirmed    = "booking_confirmed"
	EmailBookingRejected     = "booking_rejected"
	EmailHotelBookingRequest = "hotel_booking_request"
	EmailHotelBookingCancel  = "hotel_booking_cancel"
	EmailContactUsGeneral    = "contact_us_general"
	EmailContactUsBooking    = "contact_us_booking"
	EmailForgotPassword      = "forgot_password"
	EmailAccountActivated    = "account_activated"
)

const (
	ContactUsGeneral = "general"
	ContactUsBooking = "booking"
)

const (
	StatusUserActive            = "Active"
	StatusUserWaitingApproval   = "Waiting Approval"
	StatusUserInactive          = "Inactive"
	StatusUserReject            = "Reject"
	StatusUserWaitingApprovalID = 1
	StatusUserActiveID          = 2
	StatusUserRejectID          = 3
	StatusUserInactiveID        = 4
)

const (
	SupportEmail = "support_email"
)

const (
	StatusEmailSuccess   = "Success"
	StatusEmailFailed    = "Failed"
	StatusEmailPending   = "Pending"
	StatusEmailPendingID = 1
	StatusEmailSuccessID = 2
	StatusEmailFailedID  = 3
)

const (
	RoomPrice = "Room Price"
	UnitNight = "night"
	UnitPax   = "pax"
)
