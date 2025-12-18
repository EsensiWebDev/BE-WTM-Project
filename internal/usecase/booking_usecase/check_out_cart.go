package booking_usecase

import (
	"context"
	"fmt"
	"strings"
	"time"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/currency"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (bu *BookingUsecase) CheckOutCart(ctx context.Context) (*bookingdto.CheckOutCartResponse, error) {
	var invoices []entity.Invoice
	var bookingID uint

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

		bookingID = booking.ID // Capture booking ID for email function

		// 3. Update status to "in review"
		if err := bu.bookingRepo.UpdateBookingStatus(txCtx, booking.ID, constant.StatusBookingWaitingApprovalID); err != nil {
			logger.Error(ctx, "failed to update booking status", err.Error())
			return fmt.Errorf("failed to update booking status: %s", err.Error())
		}

		var countExpired int
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
					CompanyAgent:    user.AgentCompanyName,
					Agent:           user.FullName,
					Email:           user.Email,
					Hotel:           detail.DetailRooms.HotelName,
					Guest:           detail.Guest,
					CheckIn:         detail.CheckInDate.Format(time.DateOnly),
					CheckOut:        detail.CheckOutDate.Format(time.DateOnly),
					SubBookingID:    detail.SubBookingID,
					BedType:         detail.BedType,         // Selected bed type from cart
					AdditionalNotes: detail.AdditionalNotes, // Admin/agent notes from cart/booking
				},
			}
			timeNow := time.Now()
			if detail.CheckInDate.Before(timeNow) {
				countExpired++
			}

			var totalPrice float64
			var descriptionItems []entity.DescriptionInvoice
			var detailPromo entity.DetailPromo
			nights := int(detail.CheckOutDate.Sub(detail.CheckInDate).Hours() / 24)
			oriPrice := detail.RoomPrice.Price
			priceRoom := float64(nights) * oriPrice
			roomPrice := oriPrice

			// Get currency from booking detail (snapshot at booking time)
			bookingCurrency := detail.Currency
			if bookingCurrency == "" {
				bookingCurrency = user.Currency
			}
			if bookingCurrency == "" {
				bookingCurrency = "IDR" // Default fallback
			}

			if detail.Promo != nil {
				promo := detail.Promo
				detailPromo, err = bu.generateDetailPromo(promo)
				if err != nil {
					logger.Error(ctx, "failed to generate detail promo", err.Error())
				}
				// snapshot promo both on booking detail and on invoice detail
				detail.DetailPromos = detailPromo
				invoiceData.DetailInvoice.Promo = detailPromo
				nights := int(detail.CheckOutDate.Sub(detail.CheckInDate).Hours() / 24)
				switch detail.Promo.PromoTypeID {
				case constant.PromoTypeFixedPriceID:
					// Use Prices map for multi-currency support
					if len(promo.Detail.Prices) > 0 {
						// Get price for the booking currency
						if price, _, err := currency.GetPriceForCurrency(promo.Detail.Prices, bookingCurrency); err == nil {
							roomPrice = price
						} else {
							// Fallback to FixedPrice if Prices not available (backward compatibility)
							if promo.Detail.FixedPrice > 0 {
								roomPrice = promo.Detail.FixedPrice
							}
						}
					} else if promo.Detail.FixedPrice > 0 {
						// Backward compatibility: use FixedPrice if Prices not set
						roomPrice = promo.Detail.FixedPrice
					}
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
			var otherPreferences []string
			for _, additional := range detail.BookingDetailsAdditional {
				bookingDetailAdditionalName = append(bookingDetailAdditionalName, additional.NameAdditional)

				itemAdditional := entity.DescriptionInvoice{
					Description: additional.NameAdditional,
					Category:    additional.Category,
					IsRequired:  additional.IsRequired,
				}

				// Handle price-based additionals
				if additional.Category == constant.AdditionalServiceCategoryPrice && additional.Price != nil {
					quantity := 1
					price := *additional.Price
					priceAdditional := float64(quantity) * price
					itemAdditional.Quantity = quantity
					itemAdditional.Unit = constant.UnitPax
					itemAdditional.Price = price
					itemAdditional.Total = priceAdditional
					totalPrice += priceAdditional
				} else if additional.Category == constant.AdditionalServiceCategoryPax && additional.Pax != nil {
					// Handle pax-based additionals (informational only, no charge)
					itemAdditional.Quantity = *additional.Pax
					itemAdditional.Unit = constant.UnitPax
					itemAdditional.Price = 0
					itemAdditional.Total = 0
					itemAdditional.Pax = additional.Pax
				}

				descriptionItems = append(descriptionItems, itemAdditional)
			}

			// Add "Other Preferences" as informational invoice lines (no charge)
			if strings.TrimSpace(detail.OtherPreferences) != "" {
				for _, p := range strings.Split(detail.OtherPreferences, ",") {
					if name := strings.TrimSpace(p); name != "" {
						otherPreferences = append(otherPreferences, name)
						itemPref := entity.DescriptionInvoice{
							Description: name,
							Quantity:    1,
							Unit:        "preference",
							Price:       0,
							Total:       0,
						}
						descriptionItems = append(descriptionItems, itemPref)
					}
				}
			}
			//invoiceData.DetailInvoice.Promo = detail.Promo
			invoiceData.DetailInvoice.DescriptionInvoice = descriptionItems
			invoiceData.DetailInvoice.TotalPrice = totalPrice
			invoiceData.DetailInvoice.Currency = bookingCurrency // Set currency for invoice
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

		if countExpired > 0 {
			logger.Error(ctx, fmt.Sprintf("there are %d expired booking", countExpired))
			return fmt.Errorf("there are %d expired booking", countExpired)
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
		invBookingID := bookingID // Capture booking ID for goroutine
		go func(inv entity.Invoice) {
			newCtx, cancel := context.WithTimeout(context.Background(), bu.config.DurationCtxTOSlow)
			defer cancel()
			bu.sendEmailNotificationHotelConfirm(newCtx, inv.BookingDetail, invBookingID)
		}(invoice)
		resp.Invoice = append(resp.Invoice, bookingdto.DataInvoice{
			InvoiceNumber: invoice.InvoiceCode,
			DetailInvoice: invoice.DetailInvoice,
			InvoiceDate:   time.Now().Format(time.DateOnly),
		})
	}

	return resp, nil
}

func (bu *BookingUsecase) sendEmailNotificationHotelConfirm(ctx context.Context, bd entity.BookingDetail, bookingID uint) {

	logger.Info(ctx, "Data details:", bd)

	emailTemplate, err := bu.emailRepo.GetEmailTemplateByName(ctx, constant.EmailHotelBookingRequest)
	if err != nil || emailTemplate == nil {
		logger.Error(ctx, "Failed to get email template:", err)
		return
	}

	// Get full guest information from database
	guests := bu.getGuestsForEmail(ctx, bookingID)

	// Format additional services with details
	var additionalServices []AdditionalServiceEmailInfo
	for _, additional := range bd.BookingDetailsAdditional {
		serviceInfo := AdditionalServiceEmailInfo{
			Name:       additional.NameAdditional,
			Category:   additional.Category,
			IsRequired: additional.IsRequired,
		}
		if additional.Category == constant.AdditionalServiceCategoryPrice && additional.Price != nil {
			serviceInfo.Price = fmt.Sprintf("%.2f", *additional.Price)
		} else if additional.Category == constant.AdditionalServiceCategoryPax && additional.Pax != nil {
			serviceInfo.Pax = fmt.Sprintf("%d", *additional.Pax)
		}
		additionalServices = append(additionalServices, serviceInfo)
	}

	// Format bed types
	bedTypesStr := strings.Join(bd.BedTypeNames, ", ")

	data := HotelEmailData{
		Guests:             guests,
		GuestName:          bd.Guest, // Keep for backward compatibility
		Period:             fmt.Sprintf("%s to %s", bd.CheckInDate.Format("02-01-2006"), bd.CheckOutDate.Format("02-01-2006")),
		RoomType:           bd.DetailRooms.RoomTypeName,
		BedTypes:           bedTypesStr,
		Rate:               fmt.Sprintf("%.2f", bd.Price),
		BookingCode:        bd.Booking.BookingCode,
		Additional:         strings.Join(bd.BookingDetailAdditionalName, ", "), // Keep for backward compatibility
		AdditionalServices: additionalServices,
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

// GuestEmailInfo represents guest information for email template
type GuestEmailInfo struct {
	Name      string
	Honorific string
	Category  string
	Age       string // formatted as string, empty if nil
}

// AdditionalServiceEmailInfo represents additional service information for email template
type AdditionalServiceEmailInfo struct {
	Name       string
	Category   string
	Price      string // formatted as string, empty if not price-based
	Pax        string // formatted as string, empty if not pax-based
	IsRequired bool
}

// getGuestsForEmail retrieves guest information from database for email template
func (bu *BookingUsecase) getGuestsForEmail(ctx context.Context, bookingID uint) []GuestEmailInfo {
	var guests []GuestEmailInfo

	// Get full guest details from repository
	bookingGuests, err := bu.bookingRepo.GetBookingGuests(ctx, bookingID)
	if err != nil {
		logger.Error(ctx, "Failed to get booking guests for email", err.Error())
		return guests
	}

	// Convert model guests to email info
	for _, g := range bookingGuests {
		ageStr := ""
		if g.Age != nil {
			ageStr = fmt.Sprintf("%d", *g.Age)
		}
		guests = append(guests, GuestEmailInfo{
			Name:      g.Name,
			Honorific: g.Honorific,
			Category:  g.Category,
			Age:       ageStr,
		})
	}

	return guests
}

type HotelEmailData struct {
	Guests             []GuestEmailInfo
	GuestName          string // Keep for backward compatibility
	Period             string
	RoomType           string
	BedTypes           string
	Rate               string
	BookingCode        string
	Remark             string
	Additional         string // Keep for backward compatibility (comma-separated names)
	AdditionalServices []AdditionalServiceEmailInfo
	SystemSignature    string // bisa berupa teks atau <img src="...">
}
