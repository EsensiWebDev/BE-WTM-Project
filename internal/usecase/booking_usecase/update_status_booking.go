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

	notifSettings, err := bu.userRepo.GetNotificationSettings(ctx, booking.AgentID)
	if err != nil {
		logger.Error(ctx, "Failed to get notification settings:", err.Error())
		return
	}

	var templateName, title, message, typeNotif string
	var notifAllowed NotifAgentAllowed
	switch statusID {
	case constant.StatusBookingConfirmedID:
		templateName = constant.EmailBookingConfirmed
		typeNotif = constant.ConstBooking
		message = fmt.Sprintf("Your booking has been confirmed, please check Booking ID: %s", booking.BookingCode)
		title = "Booking Status Confirmed"

	case constant.StatusBookingRejectedID:
		templateName = constant.EmailBookingRejected
		typeNotif = constant.ConstReject
		message = fmt.Sprintf("Your booking has been rejected, please check Booking ID: %s", booking.BookingCode)
		title = "Booking Status Rejected"

	default:
		logger.Warn(ctx, "No email template for status:", statusID)
		return
	}

	notifAllowed = getNotifAgentAllowed(notifSettings, typeNotif)
	redirectURL := fmt.Sprintf("%s/history-booking?search_by=booking_id&search=%s", bu.config.URLFEAgent, booking.BookingCode)

	if notifAllowed.WebNotif {
		go func() {
			newCtx, cancel := context.WithTimeout(context.Background(), bu.config.DurationCtxTOSlow)
			defer cancel()
			notification := entity.Notification{
				UserID:      booking.AgentID,
				Title:       title,
				Message:     message,
				RedirectURL: redirectURL,
				Type:        typeNotif,
			}

			if err = bu.notifRepo.CreateNotification(newCtx, &notification); err != nil {
				logger.Error(ctx, "Failed to create notification:", err.Error())
				return
			}
			logger.Info(ctx, "Web notification created for agent:", booking.AgentID)
		}()

	}

	if notifAllowed.EmailNotif {
		go func() {
			newCtx, cancel := context.WithTimeout(context.Background(), bu.config.DurationCtxTOSlow)
			defer cancel()
			emailTemplate, err := bu.emailRepo.GetEmailTemplateByName(newCtx, templateName)
			if err != nil || emailTemplate == nil {
				logger.Error(newCtx, "Failed to get email template:", err)
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
				BookingLink:     redirectURL,
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

			emailTo := booking.AgentEmail

			emailLog := entity.EmailLog{
				To:              emailTo,
				Subject:         subjectParsed,
				Body:            bodyHTML,
				EmailTemplateID: uint(emailTemplate.ID),
			}
			metadataLog := entity.MetadataEmailLog{AgentName: booking.AgentName}
			emailLog.Meta = &metadataLog

			var dataEmail bool
			statusEmailID := constant.StatusEmailSuccessID
			if err = bu.emailRepo.CreateEmailLog(newCtx, &emailLog); err != nil {
				logger.Error(newCtx, "Failed to create email log:", err)
				dataEmail = false
			} else {
				dataEmail = true
			}

			err = bu.emailSender.Send(newCtx, constant.ScopeAgent, emailTo, subjectParsed, bodyHTML, "Please view this email in HTML format.")
			if err != nil {
				logger.Error(newCtx, "Failed to sending email:", err.Error())
				statusEmailID = constant.StatusEmailFailedID
				metadataLog.Notes = fmt.Sprintf("Failed to send email: %s", err.Error())
				emailLog.Meta = &metadataLog
			}

			if dataEmail {
				emailLog.StatusID = uint(statusEmailID)
				if err := bu.emailRepo.UpdateStatusEmailLog(newCtx, &emailLog); err != nil {
					logger.Error(newCtx, "Failed to update email log:", err.Error())
				}
			}
		}()

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

func getNotifAgentAllowed(notifSettings []entity.UserNotificationSetting, typeNotif string) NotifAgentAllowed {
	result := NotifAgentAllowed{}
	for _, setting := range notifSettings {
		if setting.IsEnabled && setting.Type == typeNotif {
			switch setting.Channel {
			case constant.ConstEmail:
				result.EmailNotif = true
			case constant.ConstWeb:
				result.WebNotif = true
			}
		}
	}
	return result
}

type NotifAgentAllowed struct {
	EmailNotif bool
	WebNotif   bool
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
