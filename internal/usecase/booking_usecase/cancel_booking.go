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

func (bu *BookingUsecase) CancelBooking(ctx context.Context, req *bookingdto.CancelBookingRequest) error {
	// Get agent Id from context
	userCtx, err := bu.middleware.GenerateUserFromContext(ctx)
	if err != nil {
		logger.Error(ctx, "failed to get user from context", err.Error())
		return fmt.Errorf("failed to get user from context: %s", err.Error())
	}

	if userCtx == nil {
		logger.Error(ctx, "user context is nil")
		return fmt.Errorf("user not found in context")
	}

	agentID := userCtx.ID

	bookingDetail, err := bu.bookingRepo.CancelBooking(ctx, agentID, req.SubBookingID)
	if err != nil {
		logger.Error(ctx, "failed to cancel booking", err.Error())
		return err
	}

	go func() {
		newCtx, cancel := context.WithTimeout(context.Background(), bu.config.DurationCtxTOSlow)
		defer cancel()
		bu.sendEmailNotificationHotelCancel(newCtx, bookingDetail)
	}()

	return nil
}

func (bu *BookingUsecase) sendEmailNotificationHotelCancel(ctx context.Context, bd *entity.BookingDetail) {

	if bd == nil {
		logger.Error(ctx, "Booking Detail Data is Empty")
		return
	}

	emailTemplate, err := bu.emailRepo.GetEmailTemplateByName(ctx, constant.EmailHotelBookingCancel)
	if err != nil || emailTemplate == nil {
		logger.Error(ctx, "Failed to get email template:", err)
		return
	}

	data := HotelEmailDataCancel{
		GuestName:   bd.Guest,
		Period:      fmt.Sprintf("%s to %s", bd.CheckInDate.Format("02-01-2006"), bd.CheckOutDate.Format("02-01-2006")),
		RoomType:    bd.RoomPrice.RoomType.Name,
		Rate:        fmt.Sprintf("%.2f", bd.Price),
		BookingCode: bd.Booking.BookingCode,
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
		return
	}

	err = bu.emailSender.Send(ctx, bd.RoomPrice.RoomType.Hotel.Email, subjectParsed, bodyHTML, "Please view this email in HTML format.")
	if err != nil {
		logger.Error(ctx, "Failed to send hotel booking email:", err.Error())
	}

}

type HotelEmailDataCancel struct {
	GuestName       string
	Period          string
	RoomType        string
	Rate            string
	BookingCode     string
	Remark          string
	Additional      string
	SystemSignature string // bisa berupa teks atau <img src="...">
}
