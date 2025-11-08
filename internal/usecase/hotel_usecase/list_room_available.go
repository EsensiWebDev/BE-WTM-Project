package hotel_usecase

import (
	"context"
	"fmt"
	"time"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hu *HotelUsecase) ListRoomAvailable(ctx context.Context, req *hoteldto.ListRoomAvailableRequest) (*hoteldto.ListRoomAvailableResponse, error) {
	// 1. Ambil semua room type berdasarkan hotel Id
	rooms, err := hu.hotelRepo.GetRoomTypeByHotelID(ctx, req.HotelID)
	if err != nil {
		logger.Error(ctx, "Error getting room types by hotel Id", "hotelID", req.HotelID, err.Error())
		return nil, err
	}

	// 2. Extract roomTypeIDs untuk filtering ke repo berikutnya
	roomTypeIDs := make([]uint, 0, len(rooms))
	for _, r := range rooms {
		roomTypeIDs = append(roomTypeIDs, r.ID)
	}

	monthTime, err := time.Parse("2006-01", req.Month)
	if err != nil {
		logger.Error(ctx, "Error parsing month", "month", req.Month, err.Error())
		return nil, fmt.Errorf("invalid month format: %s", err.Error())
	}

	// 3. Ambil room yang unavailable hanya berdasarkan roomTypeID
	roomUnavailable, err := hu.hotelRepo.GetRoomUnavailableByRoomTypeIDs(ctx, roomTypeIDs, monthTime)
	if err != nil {
		logger.Error(ctx, "Error getting room unavailable", "roomTypeIDs", roomTypeIDs, "month", monthTime, err.Error())
		return nil, err
	}

	// 4. Init response
	resp := &hoteldto.ListRoomAvailableResponse{}

	days, err := utils.DaysInMonth(monthTime)
	if err != nil {
		logger.Error(ctx, "Error getting days in month", "month", monthTime, err.Error())
		return nil, err
	}

	// 5. Bangun map: roomTypeID → set of tanggal unavailable
	unavailMap := make(map[uint]map[int]bool)
	for _, ru := range roomUnavailable {
		if ru.Date == nil {
			continue
		}
		day := ru.Date.Day()
		if unavailMap[ru.RoomTypeID] == nil {
			unavailMap[ru.RoomTypeID] = make(map[int]bool)
		}
		unavailMap[ru.RoomTypeID][day] = true
	}

	// 6. Construct DTO
	resp.RoomAvailable = make([]hoteldto.RoomAvailable, 0, len(rooms))
	for _, room := range rooms {
		roomAvailable := hoteldto.RoomAvailable{
			RoomTypeID:   room.ID,
			RoomTypeName: room.Name,
			Data:         make([]hoteldto.DataAvailable, 0, days),
		}

		for day := 1; day <= days; day++ {
			isUnavailable := unavailMap[room.ID][day]
			roomAvailable.Data = append(roomAvailable.Data, hoteldto.DataAvailable{
				Day:       day,
				Available: !isUnavailable,
			})
		}

		resp.RoomAvailable = append(resp.RoomAvailable, roomAvailable)
	}

	return resp, nil
}
