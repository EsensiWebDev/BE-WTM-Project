package booking_usecase

import (
	"context"
	"fmt"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (bu *BookingUsecase) UpdateStatusBooking(ctx context.Context, req *bookingdto.UpdateStatusBookingRequest) error {
	var bookingDetailIDs []uint

	if req.BookingID > 0 {
		booking, err := bu.bookingRepo.GetBookingByID(ctx, req.BookingID)
		if err != nil {
			logger.Error(ctx, "failed to get bookings", err.Error())
			return err
		}
		bookingDetailIDs = make([]uint, 0, len(booking.BookingDetails))
		for _, detail := range booking.BookingDetails {
			bookingDetailIDs = append(bookingDetailIDs, detail.ID)
		}
	} else {
		bookingDetailIDs = append(bookingDetailIDs, req.BookingDetailID)
	}

	bookingDetails, err := bu.bookingRepo.UpdateBookingDetailStatus(ctx, bookingDetailIDs, req.StatusID)
	if err != nil {
		logger.Error(ctx, "failed to update status booking", err.Error())
		return err
	}

	go func() {
		newCtx, cancel := context.WithTimeout(context.Background(), bu.config.DurationCtxTOSlow)
		defer cancel()
		bu.sendEmailNotificationHotel(newCtx, bookingDetails, req.StatusID, req.Reason)
	}()

	go func() {
		newCtx, cancel := context.WithTimeout(context.Background(), bu.config.DurationCtxTOSlow)
		defer cancel()
		bu.sendEmailNotificationAgent(newCtx, bookingDetails, req.StatusID, req.Reason)
	}()

	return nil
}

func (bu *BookingUsecase) sendEmailNotificationAgent(ctx context.Context, details []entity.BookingDetail, statusID uint, rejectionReason string) {
	if len(details) == 0 {
		logger.Warn(ctx,
			"No booking details provided for email notification")
		return
	}

	booking := details[0].Booking

	var templateName string
	switch statusID {
	case constant.StatusBookingApprovedID:
		templateName = constant.EmailBookingConfirmed
	case constant.StatusBookingRejectedID:
		templateName = constant.EmailBookingRejected
	default:
		logger.Warn(ctx,
			"No email template for status:", statusID)
		return
	}

	emailTemplate, err := bu.emailRepo.GetEmailTemplateByName(ctx, templateName)
	if err != nil || emailTemplate == nil {
		logger.Error(ctx, "Failed to get email template:", err)
		return
	}

	var subBookings []SubBookingData
	for i, bd := range details {
		subBookings = append(subBookings, SubBookingData{
			Index:     i + 1,
			HotelName: bd.DetailRooms.RoomTypeName,
			CheckIn:   bd.CheckInDate.Format("02-01-2006"),
			CheckOut:  bd.CheckOutDate.Format("02-01-2006"),
		})
	}

	data := BookingEmailData{
		AgentName:       booking.AgentName,
		BookingID:       booking.BookingCode,
		GuestName:       details[0].Guest, // atau ambil dari booking.Guests[0] jika tersedia
		BookingLink:     fmt.Sprintf("https://hotelbox.com/booking-history/%s", booking.BookingCode),
		RejectionReason: rejectionReason,
		HomePageLink:    "https://hotelbox.com",
		SubBookings:     subBookings,
	}

	subjectParsed, err := utils.ParseTemplate(emailTemplate.Subject, data)
	bodyHTML, err := utils.ParseTemplate(emailTemplate.Body, data)
	if err != nil {
		logger.Error(ctx, "Failed to parse email template:", err)
		return
	}

	err = bu.emailSender.Send(ctx, booking.AgentEmail, subjectParsed, bodyHTML, "Please view this email in HTML format.")
	if err != nil {
		logger.Error(ctx, "Failed to send booking email:", err.Error())
	}
}

func (bu *BookingUsecase) sendEmailNotificationHotel(ctx context.Context, details []entity.BookingDetail, statusID uint, rejectionReason string) {
	if len(details) == 0 {
		logger.Warn(ctx,
			"No booking details provided for hotel email notification")
		return
	}

	var templateName string
	switch statusID {
	case constant.StatusBookingApprovedID:
		templateName = constant.EmailHotelBookingRequest
	default:
		logger.Warn(ctx,
			"No email template for status:", statusID)
		return
	}

	emailTemplate, err := bu.emailRepo.GetEmailTemplateByName(ctx, templateName)
	if err != nil || emailTemplate == nil {
		logger.Error(ctx, "Failed to get email template:", err)
		return
	}

	for _, bd := range details {

		hotel, err := bu.hotelRepo.GetHotelByID(ctx, bd.RoomType.HotelID, constant.RoleAdmin)
		if err != nil || hotel == nil {
			logger.Error(ctx, "Failed to get hotel by Id:", err)
			continue
		}
		data := HotelEmailData{
			GuestName:   bd.Guest,
			Period:      fmt.Sprintf("%s to %s", bd.CheckInDate.Format("02-01-2006"), bd.CheckOutDate.Format("02-01-2006")),
			RoomType:    bd.DetailRooms.RoomTypeName,
			Rate:        fmt.Sprintf("%.2f", bd.Price),
			BookingCode: bd.Booking.BookingCode,
			Remark:      rejectionReason,
			Additional:  strings.Join(bd.BookingDetailAdditionalName, ", "),
		}

		if emailTemplate.IsSignatureImage && emailTemplate.Signature != "" {
			data.SystemSignature = bu.assignSignatureEmail(emailTemplate.Signature)
		}

		if data.SystemSignature == "" && emailTemplate.Signature != "" {
			data.SystemSignature = emailTemplate.Signature
		}

		subjectParsed, err := utils.ParseTemplate(emailTemplate.Subject, data)
		bodyHTML, err := utils.ParseTemplate(emailTemplate.Body, data)
		if err != nil {
			logger.Error(ctx, "Failed to parse hotel email template:", err)
			continue
		}

		err = bu.emailSender.Send(ctx, hotel.Email, subjectParsed, bodyHTML, "Please view this email in HTML format.")
		if err != nil {
			logger.Error(ctx, "Failed to send hotel booking email:", err.Error())
		}
	}
}

func (bu *BookingUsecase) assignSignatureEmail(emailSignature string) string {
	if emailSignature == "" {
		return ""
	}
	if strings.HasPrefix(emailSignature, "http://") || strings.HasPrefix(emailSignature, "https://") {
		return fmt.Sprintf(`<img src="%s" alt="Signature" style="width:150px;">`, emailSignature)
	}
	return emailSignature
}

type SubBookingData struct {
	Index     int
	HotelName string
	CheckIn   string // Format: "02-01-2025"
	CheckOut  string
}

type BookingEmailData struct {
	AgentName       string
	BookingID       string
	GuestName       string
	BookingLink     string
	SubBookings     []SubBookingData
	RejectionReason string // hanya dipakai untuk rejected
	HomePageLink    string // hanya dipakai untuk rejected
}

type HotelEmailData struct {
	GuestName       string
	Period          string
	RoomType        string
	Rate            string
	BookingCode     string
	Remark          string
	Additional      string
	SystemSignature string // bisa berupa teks atau <img src="...">
}
