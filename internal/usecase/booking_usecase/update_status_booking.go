package booking_usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (bu *BookingUsecase) UpdateStatusBooking(ctx context.Context, req *bookingdto.UpdateStatusRequest, scope string) error {
	var bookingDetailIDs []uint
	var err error

	var types string

	if req.BookingID != "" {
		bookingDetailIDs, err = bu.bookingRepo.GetBookingDetailIDsByBookingCode(ctx, req.BookingID)
		if err != nil {
			logger.Error(ctx, "failed to get booking detail IDs by booking code", err.Error())
			return err
		}
		types = constant.ConstBooking

	} else {
		if req.SubBookingID == "" {
			return errors.New("sub booking ID cannot be empty")
		}
		id, err := bu.bookingRepo.GetIDBySubBookingID(ctx, req.SubBookingID)
		if err != nil {
			logger.Error(ctx, "failed to get ID by sub booking ID", err.Error())
			return err
		}
		bookingDetailIDs = append(bookingDetailIDs, id)
		types = constant.ConstSubBooking

	}

	switch scope {
	case constant.ConstBooking:
		if req.StatusID == constant.StatusBookingWaitingApprovalID {
			logger.Warn(ctx, "Booking status cannot be changed to waiting approval")
			return errors.New("booking status cannot be changed to waiting approval")
		}
		bookingDetails, guests, err := bu.bookingRepo.UpdateBookingDetailStatusBooking(ctx, bookingDetailIDs, req.StatusID)
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
			bu.sendEmailNotificationAgent(newCtx, bookingDetails, req.StatusID, req.Reason, types, guests)
		}()

	case constant.ConstPayment:
		if err = bu.bookingRepo.UpdateBookingDetailStatusPayment(ctx, bookingDetailIDs, req.StatusID); err != nil {
			logger.Error(ctx, "failed to update status payment", err.Error())
			return err
		}
	}

	return nil
}

func (bu *BookingUsecase) sendEmailNotificationAgent(ctx context.Context, details []entity.BookingDetail, statusID uint, rejectionReason, types string, guests []string) {
	if len(details) == 0 {
		logger.Warn(ctx,
			"No booking details provided for email notification")
		return
	}

	booking := details[0].Booking

	var templateName string
	switch statusID {
	case constant.StatusBookingConfirmedID:
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
	for _, bd := range details {
		subBookings = append(subBookings, SubBookingData{
			SubBookingID: bd.SubBookingID,
			Guest:        bd.Guest,
			HotelName:    bd.DetailRooms.HotelName,
			CheckIn:      bd.CheckInDate.Format("02-01-2006"),
			CheckOut:     bd.CheckOutDate.Format("02-01-2006"),
		})
	}

	data := BookingEmailData{
		AgentName:       booking.AgentName,
		BookingID:       booking.BookingCode,
		GuestName:       strings.Join(guests, ", "),
		BookingLink:     fmt.Sprintf("%s/booking-history", bu.config.URLFEAgent),
		RejectionReason: rejectionReason,
		HomePageLink:    bu.config.URLFEAgent,
		SubBookings:     subBookings,
	}

	switch types {
	case constant.ConstBooking:
		data.ID = booking.BookingCode
	case constant.ConstSubBooking:
		data.ID = details[0].SubBookingID
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
	case constant.StatusBookingConfirmedID:
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

		hotel, err := bu.hotelRepo.GetHotelByID(ctx, bd.RoomPrice.RoomType.HotelID, constant.RoleAdmin)
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
	SubBookingID string
	Guest        string
	HotelName    string
	CheckIn      string // Format: "02-01-2025"
	CheckOut     string
}

type BookingEmailData struct {
	ID              string
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
