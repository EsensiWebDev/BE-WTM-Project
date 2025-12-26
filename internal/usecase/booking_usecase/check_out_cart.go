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

			// Get currency from booking detail (snapshot at booking time)
			bookingCurrency := detail.Currency
			if bookingCurrency == "" {
				bookingCurrency = user.Currency
			}
			if bookingCurrency == "" {
				bookingCurrency = "IDR" // Default fallback
			}

			// Get base price from multi-currency Prices map, fallback to deprecated Price field
			var oriPrice float64
			if len(detail.RoomPrice.Prices) > 0 {
				// Use Prices map for multi-currency support
				if price, _, err := currency.GetPriceForCurrency(detail.RoomPrice.Prices, bookingCurrency); err == nil {
					oriPrice = price
				} else {
					// Fallback to Price field if currency not found in Prices map
					oriPrice = detail.RoomPrice.Price
				}
			} else {
				// Fallback to deprecated Price field if Prices map is empty
				oriPrice = detail.RoomPrice.Price
			}

			priceRoom := float64(nights) * oriPrice
			roomPrice := oriPrice

			if detail.Promo != nil {
				promo := detail.Promo
				detailPromo, err = bu.generateDetailPromo(promo)
				if err != nil {
					logger.Error(ctx, "failed to generate detail promo", err.Error())
				}
				// snapshot promo both on booking detail and on invoice detail
				detail.DetailPromos = detailPromo
				invoiceData.DetailInvoice.Promo = detailPromo
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
			} else {
				// No promo: multiply base price by number of nights
				roomPrice = oriPrice * float64(nights)
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

	// Consolidate invoices into one invoice with all booking details
	// Group all items by sub-booking ID for clear separation
	var consolidatedItems []entity.DescriptionInvoice
	var consolidatedTotalPrice float64
	var consolidatedCurrency string
	var consolidatedSubBookingIDs []string

	// Get all guests for the booking BEFORE they are deleted
	// This must be done before deleting guests, as they are needed for emails
	allGuests := bu.getGuestsForEmail(ctx, bookingID)
	var guestNamesList []string
	for _, guest := range allGuests {
		guestStr := fmt.Sprintf("%s %s (%s", guest.Honorific, guest.Name, guest.Category)
		if guest.Category == "Child" && guest.Age != "" {
			guestStr += fmt.Sprintf(", Age: %s", guest.Age)
		}
		guestStr += ")"
		guestNamesList = append(guestNamesList, guestStr)
	}

	// Delete all guests from cart after successful checkout (now that we have them for emails)
	if err = bu.bookingRepo.DeleteAllGuestsFromBooking(ctx, bookingID); err != nil {
		logger.Error(ctx, "failed to delete guests after checkout", err.Error())
		// Don't fail if guest deletion fails, just log it
		// The checkout is already successful
	}

	// Collect all unique hotels, guests, and dates
	hotelNamesMap := make(map[string]bool)
	var allCheckInDates []string
	var allCheckOutDates []string

	// Use the first invoice as base for consolidated invoice
	if len(invoices) == 0 {
		logger.Error(ctx, "no invoices created")
		return nil, fmt.Errorf("no invoices created")
	}

	baseInvoice := invoices[0]
	consolidatedCurrency = baseInvoice.DetailInvoice.Currency

	// Collect all items from all invoices, grouped by sub-booking ID
	for _, invoice := range invoices {
		subBookingID := invoice.DetailInvoice.SubBookingID
		consolidatedSubBookingIDs = append(consolidatedSubBookingIDs, subBookingID)
		consolidatedTotalPrice += invoice.DetailInvoice.TotalPrice

		// Collect hotel names
		hotelNamesMap[invoice.DetailInvoice.Hotel] = true

		// Collect check-in and check-out dates
		if invoice.DetailInvoice.CheckIn != "" {
			allCheckInDates = append(allCheckInDates, invoice.DetailInvoice.CheckIn)
		}
		if invoice.DetailInvoice.CheckOut != "" {
			allCheckOutDates = append(allCheckOutDates, invoice.DetailInvoice.CheckOut)
		}

		// Add a separator item for each sub-booking with hotel name
		hotelName := invoice.DetailInvoice.Hotel
		separatorDescription := fmt.Sprintf("--- %s (Sub-Booking ID: %s) ---", hotelName, subBookingID)
		if len(consolidatedItems) > 0 {
			separatorItem := entity.DescriptionInvoice{
				Description: separatorDescription,
				Quantity:    0,
				Unit:        "separator",
				Price:       0,
				Total:       0,
			}
			consolidatedItems = append(consolidatedItems, separatorItem)
		} else {
			// Add header for first sub-booking
			headerItem := entity.DescriptionInvoice{
				Description: separatorDescription,
				Quantity:    0,
				Unit:        "separator",
				Price:       0,
				Total:       0,
			}
			consolidatedItems = append(consolidatedItems, headerItem)
		}

		// Add all items from this invoice
		for _, item := range invoice.DetailInvoice.DescriptionInvoice {
			consolidatedItems = append(consolidatedItems, item)
		}
	}

	// Build consolidated hotel names string
	var hotelNamesList []string
	for hotelName := range hotelNamesMap {
		hotelNamesList = append(hotelNamesList, hotelName)
	}
	consolidatedHotelNames := strings.Join(hotelNamesList, ", ")
	if consolidatedHotelNames == "" {
		consolidatedHotelNames = "Multiple Hotels"
	}

	// Build consolidated guest names string
	consolidatedGuestNames := strings.Join(guestNamesList, ", ")
	if consolidatedGuestNames == "" {
		consolidatedGuestNames = "Multiple Guests"
	}

	// Build consolidated check-in and check-out dates strings
	consolidatedCheckIn := strings.Join(allCheckInDates, ", ")
	consolidatedCheckOut := strings.Join(allCheckOutDates, ", ")

	// Create consolidated invoice detail
	consolidatedDetailInvoice := entity.DetailInvoice{
		CompanyAgent:       baseInvoice.DetailInvoice.CompanyAgent,
		Agent:              baseInvoice.DetailInvoice.Agent,
		Email:              baseInvoice.DetailInvoice.Email,
		Hotel:              consolidatedHotelNames,
		Guest:              consolidatedGuestNames,
		CheckIn:            consolidatedCheckIn,
		CheckOut:           consolidatedCheckOut,
		SubBookingID:       strings.Join(consolidatedSubBookingIDs, ", "),
		BedType:            "", // Will be shown per sub-booking in items
		AdditionalNotes:    "",
		DescriptionInvoice: consolidatedItems,
		Promo:              entity.DetailPromo{}, // Promos are per sub-booking
		TotalPrice:         consolidatedTotalPrice,
		Currency:           consolidatedCurrency,
	}

	// Create consolidated invoice response
	resp := &bookingdto.CheckOutCartResponse{
		Invoice: []bookingdto.DataInvoice{
			{
				InvoiceNumber: baseInvoice.InvoiceCode, // Use first invoice code as consolidated code
				DetailInvoice: consolidatedDetailInvoice,
				InvoiceDate:   time.Now().Format(time.DateOnly),
			},
		},
	}

	// Group bookings by hotel email for consolidated email sending
	hotelEmailMap := make(map[string][]entity.BookingDetail)
	for _, invoice := range invoices {
		hotelEmail := invoice.BookingDetail.RoomPrice.RoomType.Hotel.Email
		hotelEmailMap[hotelEmail] = append(hotelEmailMap[hotelEmail], invoice.BookingDetail)
	}

	// Send one consolidated email per hotel
	// Pass guests as parameter since they're already retrieved and will be deleted
	for hotelEmail, bookingDetails := range hotelEmailMap {
		go func(email string, details []entity.BookingDetail, guests []GuestEmailInfo) {
			newCtx, cancel := context.WithTimeout(context.Background(), bu.config.DurationCtxTOSlow)
			defer cancel()
			bu.sendConsolidatedEmailNotificationHotelConfirm(newCtx, details, bookingID, guests)
		}(hotelEmail, bookingDetails, allGuests)
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

	// Calculate room rate in IDR (always use IDR for hotel emails)
	// Note: bd.RoomPrice.Prices contains all currencies for the selected room price option.
	// For example, if agent selected room price ID 3 with {"IDR": 250000, "KRW": 10000},
	// and agent checked out using KRW, we extract IDR 250000 from the same Prices map.
	nights := int(bd.CheckOutDate.Sub(bd.CheckInDate).Hours() / 24)
	var rateIDR float64

	// Get base price in IDR from RoomPrice.Prices map
	// This gets the IDR price from the SAME room price option that the agent selected
	var basePriceIDR float64
	if len(bd.RoomPrice.Prices) > 0 {
		if price, _, err := currency.GetPriceForCurrency(bd.RoomPrice.Prices, "IDR"); err == nil {
			basePriceIDR = price
		} else {
			// Fallback to Price field if IDR not found in Prices map
			basePriceIDR = bd.RoomPrice.Price
		}
	} else {
		// Fallback to deprecated Price field if Prices map is empty
		basePriceIDR = bd.RoomPrice.Price
	}

	// Recalculate rate in IDR with promo if applicable
	if bd.Promo != nil {
		promo := bd.Promo
		switch promo.PromoTypeID {
		case constant.PromoTypeFixedPriceID:
			// Use Prices map for multi-currency support, get IDR price
			if len(promo.Detail.Prices) > 0 {
				if price, _, err := currency.GetPriceForCurrency(promo.Detail.Prices, "IDR"); err == nil {
					rateIDR = price
				} else {
					// Fallback to FixedPrice if IDR not available (backward compatibility)
					if promo.Detail.FixedPrice > 0 {
						rateIDR = promo.Detail.FixedPrice
					}
				}
			} else if promo.Detail.FixedPrice > 0 {
				// Backward compatibility: use FixedPrice if Prices not set
				rateIDR = promo.Detail.FixedPrice
			}
			if promo.Duration > nights {
				rateIDR += float64(nights-promo.Duration) * basePriceIDR
			}
		case constant.PromoTypeDiscountID:
			rateIDR = (100 - promo.Detail.DiscountPercentage) / 100 * basePriceIDR * float64(nights)
			if promo.Duration > nights {
				rateIDR += float64(nights-promo.Duration) * basePriceIDR
			}
		default:
			rateIDR = basePriceIDR * float64(nights)
		}
	} else {
		// No promo: multiply base price by number of nights
		rateIDR = basePriceIDR * float64(nights)
	}

	// Fetch original RoomTypeAdditional to get IDR prices for additional services
	var roomTypeAdditionalIDs []uint
	for _, additional := range bd.BookingDetailsAdditional {
		roomTypeAdditionalIDs = append(roomTypeAdditionalIDs, additional.RoomTypeAdditionalID)
	}

	// Map of RoomTypeAdditionalID to RoomTypeAdditional for quick lookup
	roomTypeAdditionalsMap := make(map[uint]entity.RoomTypeAdditional)
	if len(roomTypeAdditionalIDs) > 0 {
		roomTypeAdditionals, err := bu.hotelRepo.GetRoomTypeAdditionalsByIDs(ctx, roomTypeAdditionalIDs)
		if err != nil {
			logger.Error(ctx, "Failed to get room type additionals for IDR prices", err.Error())
		} else {
			for _, rta := range roomTypeAdditionals {
				roomTypeAdditionalsMap[rta.ID] = rta
			}
		}
	}

	// Format additional services with IDR prices
	var additionalServices []AdditionalServiceEmailInfo
	for _, additional := range bd.BookingDetailsAdditional {
		serviceInfo := AdditionalServiceEmailInfo{
			Name:       additional.NameAdditional,
			Category:   additional.Category,
			IsRequired: additional.IsRequired,
		}
		if additional.Category == constant.AdditionalServiceCategoryPrice {
			// Get IDR price from original RoomTypeAdditional
			if rta, exists := roomTypeAdditionalsMap[additional.RoomTypeAdditionalID]; exists {
				if len(rta.Prices) > 0 {
					if priceIDR, _, err := currency.GetPriceForCurrency(rta.Prices, "IDR"); err == nil {
						serviceInfo.Price = fmt.Sprintf("%.2f", priceIDR)
					} else if rta.Price != nil {
						// Fallback to Price field if IDR not found
						serviceInfo.Price = fmt.Sprintf("%.2f", *rta.Price)
					}
				} else if rta.Price != nil {
					// Fallback to Price field if Prices map is empty
					serviceInfo.Price = fmt.Sprintf("%.2f", *rta.Price)
				}
			} else if additional.Price != nil {
				// Last resort: use the stored price (may not be IDR, but better than nothing)
				logger.Warn(ctx, fmt.Sprintf("RoomTypeAdditional not found for ID %d, using stored price", additional.RoomTypeAdditionalID))
				serviceInfo.Price = fmt.Sprintf("%.2f", *additional.Price)
			}
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
		Rate:               currency.FormatCurrency(rateIDR, "IDR", "IDR"), // Use IDR price with currency formatting
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
		HotelName:   bd.RoomPrice.RoomType.Hotel.Name,
		BookingCode: bd.Booking.BookingCode,
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

// sendConsolidatedEmailNotificationHotelConfirm sends one email per hotel with all booking details
func (bu *BookingUsecase) sendConsolidatedEmailNotificationHotelConfirm(ctx context.Context, bookingDetails []entity.BookingDetail, bookingID uint, guests []GuestEmailInfo) {
	if len(bookingDetails) == 0 {
		logger.Warn(ctx, "No booking details provided for consolidated email")
		return
	}

	logger.Info(ctx, fmt.Sprintf("Sending consolidated email for %d booking details", len(bookingDetails)))

	emailTemplate, err := bu.emailRepo.GetEmailTemplateByName(ctx, constant.EmailHotelBookingRequest)
	if err != nil || emailTemplate == nil {
		logger.Error(ctx, "Failed to get email template:", err)
		return
	}

	// Use guests passed as parameter (retrieved before deletion)
	// This ensures all guests are included even after they're deleted from the database

	// Get hotel info from first booking detail (all should be same hotel)
	firstBD := bookingDetails[0]
	hotelEmail := firstBD.RoomPrice.RoomType.Hotel.Email
	hotelName := firstBD.RoomPrice.RoomType.Hotel.Name
	bookingCode := firstBD.Booking.BookingCode

	// Build consolidated booking details
	var consolidatedBookings []ConsolidatedBookingDetail
	for index, bd := range bookingDetails {
		// Calculate room rate in IDR (always use IDR for hotel emails)
		nights := int(bd.CheckOutDate.Sub(bd.CheckInDate).Hours() / 24)
		var rateIDR float64

		// Get base price in IDR from RoomPrice.Prices map
		var basePriceIDR float64
		if len(bd.RoomPrice.Prices) > 0 {
			if price, _, err := currency.GetPriceForCurrency(bd.RoomPrice.Prices, "IDR"); err == nil {
				basePriceIDR = price
			} else {
				basePriceIDR = bd.RoomPrice.Price
			}
		} else {
			basePriceIDR = bd.RoomPrice.Price
		}

		// Recalculate rate in IDR with promo if applicable
		if bd.Promo != nil {
			promo := bd.Promo
			switch promo.PromoTypeID {
			case constant.PromoTypeFixedPriceID:
				if len(promo.Detail.Prices) > 0 {
					if price, _, err := currency.GetPriceForCurrency(promo.Detail.Prices, "IDR"); err == nil {
						rateIDR = price
					} else if promo.Detail.FixedPrice > 0 {
						rateIDR = promo.Detail.FixedPrice
					}
				} else if promo.Detail.FixedPrice > 0 {
					rateIDR = promo.Detail.FixedPrice
				}
				if promo.Duration > nights {
					rateIDR += float64(nights-promo.Duration) * basePriceIDR
				}
			case constant.PromoTypeDiscountID:
				rateIDR = (100 - promo.Detail.DiscountPercentage) / 100 * basePriceIDR * float64(nights)
				if promo.Duration > nights {
					rateIDR += float64(nights-promo.Duration) * basePriceIDR
				}
			default:
				rateIDR = basePriceIDR * float64(nights)
			}
		} else {
			rateIDR = basePriceIDR * float64(nights)
		}

		// Fetch original RoomTypeAdditional to get IDR prices for additional services
		var roomTypeAdditionalIDs []uint
		for _, additional := range bd.BookingDetailsAdditional {
			roomTypeAdditionalIDs = append(roomTypeAdditionalIDs, additional.RoomTypeAdditionalID)
		}

		roomTypeAdditionalsMap := make(map[uint]entity.RoomTypeAdditional)
		if len(roomTypeAdditionalIDs) > 0 {
			roomTypeAdditionals, err := bu.hotelRepo.GetRoomTypeAdditionalsByIDs(ctx, roomTypeAdditionalIDs)
			if err != nil {
				logger.Error(ctx, "Failed to get room type additionals for IDR prices", err.Error())
			} else {
				for _, rta := range roomTypeAdditionals {
					roomTypeAdditionalsMap[rta.ID] = rta
				}
			}
		}

		// Format additional services with IDR prices
		var additionalServices []AdditionalServiceEmailInfo
		for _, additional := range bd.BookingDetailsAdditional {
			serviceInfo := AdditionalServiceEmailInfo{
				Name:       additional.NameAdditional,
				Category:   additional.Category,
				IsRequired: additional.IsRequired,
			}
			if additional.Category == constant.AdditionalServiceCategoryPrice {
				if rta, exists := roomTypeAdditionalsMap[additional.RoomTypeAdditionalID]; exists {
					if len(rta.Prices) > 0 {
						if priceIDR, _, err := currency.GetPriceForCurrency(rta.Prices, "IDR"); err == nil {
							serviceInfo.Price = fmt.Sprintf("%.2f", priceIDR)
						} else if rta.Price != nil {
							serviceInfo.Price = fmt.Sprintf("%.2f", *rta.Price)
						}
					} else if rta.Price != nil {
						serviceInfo.Price = fmt.Sprintf("%.2f", *rta.Price)
					}
				} else if additional.Price != nil {
					logger.Warn(ctx, fmt.Sprintf("RoomTypeAdditional not found for ID %d, using stored price", additional.RoomTypeAdditionalID))
					serviceInfo.Price = fmt.Sprintf("%.2f", *additional.Price)
				}
			} else if additional.Category == constant.AdditionalServiceCategoryPax && additional.Pax != nil {
				serviceInfo.Pax = fmt.Sprintf("%d", *additional.Pax)
			}
			additionalServices = append(additionalServices, serviceInfo)
		}

		// Format bed types
		bedTypesStr := strings.Join(bd.BedTypeNames, ", ")

		consolidatedBooking := ConsolidatedBookingDetail{
			BookingNumber:      index + 1, // 1-based numbering
			SubBookingID:       bd.SubBookingID,
			GuestName:          bd.Guest,
			Period:             fmt.Sprintf("%s to %s", bd.CheckInDate.Format("02-01-2006"), bd.CheckOutDate.Format("02-01-2006")),
			RoomType:           bd.DetailRooms.RoomTypeName,
			BedTypes:           bedTypesStr,
			Rate:               currency.FormatCurrency(rateIDR, "IDR", "IDR"), // Use IDR price with currency formatting
			AdditionalServices: additionalServices,
			Additional:         strings.Join(bd.BookingDetailAdditionalName, ", "),
		}
		consolidatedBookings = append(consolidatedBookings, consolidatedBooking)
	}

	// Build email data with consolidated bookings
	data := HotelEmailData{
		Guests:             guests,
		GuestName:          bookingDetails[0].Guest, // Keep for backward compatibility (first guest)
		Period:             "",                      // Will be shown per booking
		RoomType:           "",                      // Will be shown per booking
		BedTypes:           "",                      // Will be shown per booking
		Rate:               "",                      // Will be shown per booking
		BookingCode:        bookingCode,
		Additional:         "",                             // Will be shown per booking
		AdditionalServices: []AdditionalServiceEmailInfo{}, // Will be shown per booking
		BookingDetails:     consolidatedBookings,
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

	emailLog := entity.EmailLog{
		To:              hotelEmail,
		Subject:         subjectParsed,
		Body:            bodyHTML,
		EmailTemplateID: uint(emailTemplate.ID),
	}
	metadataLog := entity.MetadataEmailLog{
		HotelName:   hotelName,
		BookingCode: bookingCode,
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

	err = bu.emailSender.Send(ctx, constant.ScopeHotel, hotelEmail, subjectParsed, bodyHTML, "Please view this email in HTML format.")
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
	// Consolidated booking data
	BookingDetails []ConsolidatedBookingDetail // For multiple bookings in one email
}

// ConsolidatedBookingDetail represents a single booking detail for consolidated email
type ConsolidatedBookingDetail struct {
	BookingNumber      int // 1-based booking number for display
	SubBookingID       string
	GuestName          string
	Period             string
	RoomType           string
	BedTypes           string
	Rate               string
	AdditionalServices []AdditionalServiceEmailInfo
	Additional         string
}
