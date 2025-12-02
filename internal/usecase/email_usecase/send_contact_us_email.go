package email_usecase

import (
	"context"
	"fmt"
	"strings"
	"time"
	"wtm-backend/internal/dto/emaildto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (eu *EmailUsecase) SendContactUsEmail(ctx context.Context, req *emaildto.SendContactUsEmailRequest) error {
	// Tentukan nama template berdasarkan type
	var templateName string
	switch req.Type {
	case constant.ContactUsGeneral:
		templateName = constant.EmailContactUsGeneral
	case constant.ContactUsBooking:
		templateName = constant.EmailContactUsBooking
	default:
		logger.Error(ctx, "invalid contact us type", req.Type)
		return fmt.Errorf("invalid contact us type: %s", req.Type)
	}

	// Ambil template
	emailTemplate, err := eu.emailRepo.GetEmailTemplateByName(ctx, templateName)
	if err != nil {
		logger.Error(ctx, "get email template by name fail", err.Error())
		return err
	}
	if emailTemplate == nil {
		logger.Error(ctx, "email template not found", templateName)
		return fmt.Errorf("email template not found: %s", templateName)
	}

	// Render berdasar tipe
	switch req.Type {
	case constant.ContactUsGeneral:
		// Data untuk general
		data := GeneralContactEmailData{
			UserName:    req.Name,
			UserEmail:   req.Email,
			Subject:     req.Subject,
			UserMessage: req.Message,
		}

		subject, err := utils.ParseTemplate(emailTemplate.Subject, data)
		if err != nil {
			logger.Error(ctx, "parse subject general fail", err.Error())
			return err
		}
		bodyHTML, err := utils.ParseTemplate(emailTemplate.Body, data)
		if err != nil {
			logger.Error(ctx, "parse body general fail", err.Error())
			return err
		}
		bodyText := "Please view this email in HTML format."

		if err := eu.emailSender.Send(ctx, constant.SupportEmail, subject, bodyHTML, bodyText); err != nil {
			logger.Error(ctx, "send general contact email fail", err.Error())
			return err
		}
		return nil

	case constant.ContactUsBooking:
		// Validasi basic
		if req.BookingCode == "" {
			return fmt.Errorf("booking_code is required for booking contact type")
		}

		// Ambil booking utama
		booking, err := eu.bookingRepo.GetBookingByCode(ctx, req.BookingCode)
		if err != nil {
			logger.Error(ctx, "get booking by code fail", err.Error())
			return err
		}
		if booking == nil {
			return fmt.Errorf("booking not found: %s", req.BookingCode)
		}

		// Bangun daftar sub-booking
		var subBookings []SubBooking
		var guests []string

		if req.SubBookingCode != "" {
			sb, err := eu.bookingRepo.GetSubBookingByCode(ctx, req.SubBookingCode)
			if err != nil {
				logger.Error(ctx, "get sub booking by code fail", err.Error())
				return err
			}
			if sb == nil {
				return fmt.Errorf("sub booking not found: %s", req.SubBookingCode)
			}
			subBookings = append(subBookings, SubBooking{
				SubBookingCode: sb.SubBookingID,
				GuestName:      sb.Guest,
				Hotel:          safeHotelName(sb.DetailRooms.HotelName),
				CheckIn:        formatDate(sb.CheckInDate),
				CheckOut:       formatDate(sb.CheckOutDate),
			})

			// Ambil guest
			guests = append(guests, sb.Guest)
		} else {
			for _, sb := range booking.BookingDetails {
				subBookings = append(subBookings, SubBooking{
					SubBookingCode: sb.SubBookingID,
					GuestName:      sb.Guest,
					Hotel:          safeHotelName(sb.DetailRooms.HotelName),
					CheckIn:        formatDate(sb.CheckInDate),
					CheckOut:       formatDate(sb.CheckOutDate),
				})
				guests = append(guests, sb.Guest)
			}
		}

		// Data untuk template booking
		data := BookingContactEmailData{
			BookingID:    booking.BookingCode, // atau booking.BookingID sesuai field di entity
			GuestName:    strings.Join(guests, ", "),
			AgentName:    req.Name,
			AgencyName:   booking.AgentCompanyName,
			AgentEmail:   req.Email,
			AgentPhone:   booking.AgentPhoneNumber,
			AgentMessage: req.Message,
			SubBookings:  subBookings,
			// AgentPhone & AgencyName bisa diisi kalau accessible dari booking/agent
		}

		subject, err := utils.ParseTemplate(emailTemplate.Subject, data)
		if err != nil {
			logger.Error(ctx, "parse subject booking fail", err.Error())
			return err
		}
		bodyHTML, err := utils.ParseTemplate(emailTemplate.Body, data)
		if err != nil {
			logger.Error(ctx, "parse body booking fail", err.Error())
			return err
		}
		bodyText := "Please view this email in HTML format."

		if err := eu.emailSender.Send(ctx, constant.SupportEmail, subject, bodyHTML, bodyText); err != nil {
			logger.Error(ctx, "send booking contact email fail", err.Error())
			return err
		}
		return nil
	}

	return nil
}

// Normalisasi nama hotel kalau kemungkinan null/empty dari DB.
func safeHotelName(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

// Format tanggal fleksibel, sesuaikan tipe aslinya.
func formatDate(t time.Time) string {
	// contoh: 02 Jan 2006
	return t.Format("02 Jan 2006")
}

type GeneralContactEmailData struct {
	UserName    string
	UserEmail   string
	Subject     string // subject dari user (bukan subject template)
	UserMessage string
}

type SubBooking struct {
	SubBookingCode string
	GuestName      string
	Hotel          string
	CheckIn        string // bisa pakai time.Time kalau mau parsing tanggal
	CheckOut       string
}

type BookingContactEmailData struct {
	BookingID    string
	AgentName    string
	AgentEmail   string
	AgencyName   string
	AgentPhone   string
	GuestName    string
	SubBookings  []SubBooking
	AgentMessage string
}
