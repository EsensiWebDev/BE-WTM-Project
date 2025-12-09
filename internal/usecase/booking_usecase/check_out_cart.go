package booking_usecase

import (
	"context"
	"fmt"
	"strings"
	"time"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (bu *BookingUsecase) CheckOutCart(ctx context.Context) (*bookingdto.CheckOutCartResponse, error) {
	var invoices []entity.Invoice

	err := bu.dbTrx.WithTransaction(ctx, func(txCtx context.Context) error {
		// Get agent Id from context
		userCtx, err := bu.middleware.GenerateUserFromContext(txCtx)
		if err != nil {
			logger.Error(ctx, "failed to get user from context", err.Error())
			return fmt.Errorf("failed to get user from context: %s", err.Error())
		}

		if userCtx == nil {
			logger.Error(ctx, "user context is nil")
			return fmt.Errorf("user context is nil")
		}

		agentID := userCtx.ID

		user, err := bu.userRepo.GetUserByID(txCtx, agentID)
		if err != nil {
			logger.Error(ctx, "failed to get user", err.Error())
			return fmt.Errorf("failed to get user: %s", err.Error())
		}

		booking, err := bu.bookingRepo.GetCartBooking(txCtx, agentID)
		if err != nil {
			logger.Error(ctx, "failed to get card booking", err.Error())
			return fmt.Errorf("failed to get cart: %s", err.Error())
		}
		if booking == nil {
			logger.Error(ctx, "card booking is nil")
			return fmt.Errorf("no cart found")
		}

		// 3. Update status to "in review"
		if err := bu.bookingRepo.UpdateBookingStatus(txCtx, booking.ID, constant.StatusBookingWaitingApprovalID); err != nil {
			logger.Error(ctx, "failed to update booking status", err.Error())
			return fmt.Errorf("failed to update booking status: %s", err.Error())
		}

		// 4. Create Invoice Data
		for _, detail := range booking.BookingDetails {
			cancellationDate := detail.CheckInDate.AddDate(0, 0, detail.RoomPrice.RoomType.Hotel.CancellationPeriod)
			detailRoom := entity.DetailRoom{
				HotelName:     detail.RoomPrice.RoomType.Hotel.Name,
				RoomTypeName:  detail.RoomPrice.RoomType.Name,
				Capacity:      detail.RoomPrice.RoomType.MaxOccupancy,
				IsAPI:         detail.RoomPrice.RoomType.Hotel.IsAPI,
				CancelledDate: cancellationDate.Format(time.DateOnly),
			}
			detail.DetailRooms = detailRoom

			invoiceCode, err := bu.bookingRepo.GenerateCode(ctx, "invoice_codes", "INV")
			if err != nil {
				logger.Error(ctx, "failed to generate invoice code", "error", err)
				return err
			}

			if strings.TrimSpace(detail.Guest) == "" {
				logger.Error(ctx, fmt.Sprintf("guest in room %s is empty", detail.DetailRooms.RoomTypeName))
				return fmt.Errorf("guest in room %s is empty", detail.DetailRooms.RoomTypeName)
			}
			invoiceData := entity.Invoice{
				BookingDetailID: detail.ID,
				InvoiceCode:     invoiceCode,
				DetailInvoice: entity.DetailInvoice{
					CompanyAgent: user.AgentCompanyName,
					Agent:        user.FullName,
					Email:        user.Email,
					Hotel:        detail.DetailRooms.HotelName,
					Guest:        detail.Guest,
					CheckIn:      detail.CheckInDate.Format(time.DateOnly),
					CheckOut:     detail.CheckOutDate.Format(time.DateOnly),
					SubBookingID: detail.SubBookingID,
				},
			}
			var totalPrice float64
			var descriptionItems []entity.DescriptionInvoice
			var detailPromo entity.DetailPromo
			nights := int(detail.CheckOutDate.Sub(detail.CheckInDate).Hours() / 24)
			oriPrice := detail.RoomPrice.Price
			priceRoom := float64(nights) * oriPrice
			roomPrice := oriPrice

			if detail.Promo != nil {
				promo := detail.Promo
				detailPromo, err = bu.generateDetailPromo(promo)
				if err != nil {
					logger.Error(ctx, "failed to generate detail promo", err.Error())
				}
				detail.DetailPromos = detailPromo
				nights := int(detail.CheckOutDate.Sub(detail.CheckInDate).Hours() / 24)
				switch detail.Promo.PromoTypeID {
				case constant.PromoTypeFixedPriceID:
					roomPrice = promo.Detail.FixedPrice
					if promo.Duration > nights {
						roomPrice += float64(nights-promo.Duration) * oriPrice
					}
				case constant.PromoTypeDiscountID:
					roomPrice = (100 - promo.Detail.DiscountPercentage) / 100 * oriPrice * float64(nights)
					if promo.Duration > nights {
						roomPrice += float64(nights-promo.Duration) * oriPrice
					}
				default:
					roomPrice = roomPrice * float64(nights)
				}
			}
			detail.Price = roomPrice
			itemRoom := entity.DescriptionInvoice{
				Description:      detail.DetailRooms.RoomTypeName,
				Quantity:         nights,
				Unit:             constant.UnitNight,
				Price:            oriPrice,
				TotalBeforePromo: priceRoom,
				Total:            detail.Price,
			}
			totalPrice += itemRoom.Total
			descriptionItems = append(descriptionItems, itemRoom)
			var bookingDetailAdditionalName []string
			for _, additional := range detail.BookingDetailsAdditional {
				quantity := 1
				price := additional.Price
				priceAdditional := float64(quantity) * price
				bookingDetailAdditionalName = append(bookingDetailAdditionalName, additional.NameAdditional)
				itemAdditional := entity.DescriptionInvoice{
					Description: additional.NameAdditional,
					Quantity:    quantity,
					Unit:        constant.UnitPax,
					Price:       additional.Price,
					Total:       priceAdditional,
				}
				totalPrice += priceAdditional
				descriptionItems = append(descriptionItems, itemAdditional)
			}
			//invoiceData.DetailInvoice.Promo = detail.Promo
			invoiceData.DetailInvoice.DescriptionInvoice = descriptionItems
			invoiceData.DetailInvoice.TotalPrice = totalPrice
			invoiceData.BookingDetail = detail
			invoiceData.BookingDetail.Price = detail.Price
			invoiceData.BookingDetail.Booking.BookingCode = booking.BookingCode
			invoiceData.BookingDetail.BookingDetailAdditionalName = bookingDetailAdditionalName
			invoices = append(invoices, invoiceData)

			//Update Detail Booking Detail
			if err = bu.bookingRepo.UpdateDetailBookingDetail(txCtx, detail.ID, &detailRoom, &detailPromo, detail.Price, detail.BookingDetailsAdditional); err != nil {
				logger.Error(ctx, "failed to update booking", err.Error())
				return fmt.Errorf("failed to update booking: %s", err.Error())
			}
		}

		// Create Invoice
		if err = bu.bookingRepo.CreateInvoice(txCtx, invoices); err != nil {
			logger.Error(ctx, "failed to create invoice", err.Error())
			return fmt.Errorf("failed to create invoice: %s", err.Error())
		}

		return nil
	})

	if err != nil {
		logger.Error(ctx, "transaction failed in check out cart", err.Error())
		return nil, err
	}

	// sekarang invoices sudah terisi
	resp := &bookingdto.CheckOutCartResponse{
		Invoice: make([]bookingdto.DataInvoice, 0, len(invoices)),
	}
	for _, invoice := range invoices {
		go func() {
			newCtx, cancel := context.WithTimeout(context.Background(), bu.config.DurationCtxTOSlow)
			defer cancel()
			bu.sendEmailNotificationHotelConfirm(newCtx, invoice.BookingDetail)
		}()
		resp.Invoice = append(resp.Invoice, bookingdto.DataInvoice{
			InvoiceNumber: invoice.InvoiceCode,
			DetailInvoice: invoice.DetailInvoice,
			InvoiceDate:   time.Now().Format(time.DateOnly),
		})
	}

	return resp, nil
}

func (bu *BookingUsecase) sendEmailNotificationHotelConfirm(ctx context.Context, bd entity.BookingDetail) {

	logger.Info(ctx, "Data details:", bd)

	emailTemplate, err := bu.emailRepo.GetEmailTemplateByName(ctx, constant.EmailHotelBookingRequest)
	if err != nil || emailTemplate == nil {
		logger.Error(ctx, "Failed to get email template:", err)
		return
	}

	data := HotelEmailData{
		GuestName:   bd.Guest,
		Period:      fmt.Sprintf("%s to %s", bd.CheckInDate.Format("02-01-2006"), bd.CheckOutDate.Format("02-01-2006")),
		RoomType:    bd.DetailRooms.RoomTypeName,
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

	emailTo := bd.RoomPrice.RoomType.Hotel.Email

	emailLog := entity.EmailLog{
		To:              emailTo,
		Subject:         subjectParsed,
		Body:            bodyHTML,
		EmailTemplateID: uint(emailTemplate.ID),
	}
	metadataLog := entity.MetadataEmailLog{
		HotelName: bd.RoomPrice.RoomType.Hotel.Name,
	}
	emailLog.Meta = &metadataLog

	var dataEmail bool
	statusEmailID := constant.StatusEmailSuccessID
	if err = bu.emailRepo.CreateEmailLog(ctx, &emailLog); err != nil {
		logger.Error(ctx, "Failed to create email log:", err)
		dataEmail = false
	} else {
		dataEmail = true
	}

	err = bu.emailSender.Send(ctx, constant.ScopeHotel, emailTo, subjectParsed, bodyHTML, "Please view this email in HTML format.")
	if err != nil {
		logger.Error(ctx, "Failed to sending email:", err.Error())
		statusEmailID = constant.StatusEmailFailedID
		metadataLog.Notes = fmt.Sprintf("Failed to send email: %s", err.Error())
		emailLog.Meta = &metadataLog
	}

	if dataEmail {
		emailLog.StatusID = uint(statusEmailID)
		if err := bu.emailRepo.UpdateStatusEmailLog(ctx, &emailLog); err != nil {
			logger.Error(ctx, "Failed to update email log:", err.Error())
		}
	}

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
