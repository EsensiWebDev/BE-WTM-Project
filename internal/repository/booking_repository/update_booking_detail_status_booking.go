package booking_repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (br *BookingRepository) UpdateBookingDetailStatusBooking(ctx context.Context, bookingDetailIDs []uint, statusID uint) ([]entity.BookingDetail, []string, error) {
	db := br.db.GetTx(ctx)

	var updatedDetails []model.BookingDetail
	var guests []string
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Step 1: Update BookingDetail yang masih waiting approval
		res := tx.Model(&model.BookingDetail{}).
			Where("id IN ?", bookingDetailIDs).
			Where("status_booking_id = ?", constant.StatusBookingWaitingApprovalID)

		if statusID == constant.StatusBookingConfirmedID {
			res.Updates(map[string]interface{}{
				"status_booking_id": statusID,
				"approved_at":       time.Now(),
			})
		} else {
			res.Updates(map[string]interface{}{
				"status_booking_id": statusID,
			})
		}

		if err := res.Error; err != nil {
			logger.Error(ctx, "failed to update booking detail status", err.Error())
			return err
		}

		if res.RowsAffected == 0 {
			logger.Error(ctx, "no booking detail updated")
			return errors.New("no booking detail updated")
		}

		// Step 2: Ambil bookingID (cukup satu, karena pasti sama)
		var bookingID uint
		if err := tx.Model(&model.BookingDetail{}).
			Select("booking_id").
			Where("id = ?", bookingDetailIDs[0]).
			Scan(&bookingID).Error; err != nil {
			logger.Error(ctx, "failed to get booking id")
			return err
		}

		// Step 3: Cek apakah masih ada waiting approval di booking ini
		var waitingCount int64
		if err := tx.Model(&model.BookingDetail{}).
			Where("booking_id = ? AND status_booking_id = ?", bookingID, constant.StatusBookingWaitingApprovalID).
			Count(&waitingCount).Error; err != nil {
			logger.Error(ctx, "failed to count waiting scheduler status", err.Error())
			return err
		}

		if waitingCount == 0 {
			orderExpr := fmt.Sprintf(`
				CASE 
					WHEN status_booking_id = %d THEN 1
					WHEN status_booking_id = %d THEN 2
					WHEN status_booking_id = %d THEN 3
					ELSE 4
				END`,
				constant.StatusBookingRejectedID,
				constant.StatusBookingConfirmedID,
				constant.StatusBookingCanceledID,
			)
			// Step 3a: Tentukan status prioritas dari semua detail
			var detailStatus uint
			if err := tx.Model(&model.BookingDetail{}).
				Select("status_booking_id").
				Where("booking_id = ?", bookingID).
				Order(orderExpr).
				Limit(1).
				Scan(&detailStatus).Error; err != nil {
				logger.Error(ctx, "failed to get booking status", err.Error())
				return err
			}

			// Step 3b: Update Booking sesuai prioritas
			if err := tx.Model(&model.Booking{}).
				Where("id = ?", bookingID).
				Updates(map[string]interface{}{
					"status_booking_id": detailStatus,
					"approved_at":       time.Now(),
				}).Error; err != nil {
				logger.Error(ctx, "failed to update booking with priority", err.Error())
				return err
			}
		}

		// Step 4: Fetch updated BookingDetails dengan preload
		if err := tx.Preload("Booking").
			Preload("Booking.Agent").
			Preload("RoomPrice").
			Preload("RoomPrice.RoomType").
			Preload("BookingDetailsAdditional").
			Where("booking_id = ?", bookingID).
			Find(&updatedDetails).Error; err != nil {
			logger.Error(ctx, "failed to fetch updated booking details", err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		logger.Error(ctx, "failed to update status booking", err.Error())
		return nil, guests, err
	}

	var bdEntity, result []entity.BookingDetail
	if err := utils.CopyStrict(&bdEntity, &updatedDetails); err != nil {
		logger.Error(ctx, "failed to copy updated booking details", err.Error())
		return nil, guests, err
	}

	for i, detail := range updatedDetails {
		guests = append(guests, detail.Guest)
		for _, d := range bookingDetailIDs {
			if detail.ID == d {
				var dataResult entity.BookingDetail
				dataResult = bdEntity[i]
				var detailRoom entity.DetailRoom
				if err := json.Unmarshal(detail.DetailRoom, &detailRoom); err != nil {
					logger.Error(ctx, "Error marshalling room detail to JSON", err.Error())
				}
				dataResult.DetailRooms = detailRoom

				if len(detail.BookingDetailsAdditional) > 0 {
					dataResult.BookingDetailAdditionalName = make([]string, 0, len(detail.BookingDetailsAdditional))
					for _, additional := range detail.BookingDetailsAdditional {
						dataResult.BookingDetailAdditionalName = append(dataResult.BookingDetailAdditionalName, additional.NameAdditional)
					}
				}
				if detail.Booking.ID != 0 {
					dataResult.Booking.AgentName = detail.Booking.Agent.FullName
					dataResult.Booking.AgentEmail = detail.Booking.Agent.Email
				}
				result = append(result, dataResult)
				break
			}
		}
	}

	return result, guests, nil
}
